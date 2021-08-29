package main

import (
	"strings"
	"errors"
	"io"
	"fmt"
)

type Rule struct {
	Name string
	Char rune
	Action Action
	WildCard bool
}

type Action func(ctx *ParserContext, char rune) error

type Parser map[string][]Rule

type ParserContext struct {
	Stack []string
	Parser Parser
	Stream *strings.Reader
	Next bool
	Captured []*strings.Builder
	Processed bool
	Ast *Node
}


func Push(next_expr string) Action {
	return func(ctx *ParserContext, char rune) error {
		var curr_capture strings.Builder 
		ctx.Captured = append(ctx.Captured, &curr_capture)

		//fmt.Printf("PUSH: %s\n", next_expr)

		ctx.Stack = append(ctx.Stack, next_expr)

		ctx.Next = true
		ctx.Processed = true

		child := &Node{
			Type: next_expr,
			Childs: make([]*Node, 0),
			Parent: ctx.Ast,
		}

		ctx.Ast.Childs = append(ctx.Ast.Childs, child)
		ctx.Ast = child

		return nil
	}
}

func Read() Action {
	return func(ctx *ParserContext, char rune) error {
		curr_capture := ctx.Captured[len(ctx.Captured)-1]

		//fmt.Printf("READ: [%s]\n", string(char))
		curr_capture.WriteRune(char)
		ctx.Processed = true

		return nil
	}
}

func Goto(next_expr string) Action {
	return func(ctx *ParserContext, char rune) error {
		//fmt.Printf("GOTO: %s\n", next_expr)

		ctx.Next = true
		ctx.Processed = true

		ctx.Stack = append(ctx.Stack, next_expr)

		return nil
	}
}

func Return() Action {
	return func(ctx *ParserContext, char rune) error {
		curr_capture := ctx.Captured[len(ctx.Captured)-1]

		curr_capture.WriteRune(char)

		//fmt.Printf("RETURN: [%s]\n", curr_capture.String())

		ctx.Next = true
		ctx.Stack = ctx.Stack[:len(ctx.Stack)-1]

		ctx.Processed = true
		return nil
	}
}

func Pop() Action {
	return func(ctx *ParserContext, char rune) error {
		curr_capture := ctx.Captured[len(ctx.Captured)-1]

		//fmt.Printf("[%s => %s]\n", ctx.Stack[len(ctx.Stack)-1], curr_capture.String())

		ctx.Next = true
		ctx.Stack = ctx.Stack[:len(ctx.Stack)-1]
		ctx.Captured = ctx.Captured[:len(ctx.Captured)-1]

		ctx.Processed = true

		ctx.Ast.Value = curr_capture.String()
		ctx.Ast = ctx.Ast.Parent

		return nil
	}
}

func Ignore(reg string) Action {
	return func(ctx *ParserContext, char rune) error {
		if (!strings.ContainsRune(reg, char)) {
			ctx.Stream.UnreadRune()
			return nil
		}

		ctx.Processed = true

		if ctx.Stream.Len() == 0 {
			return errors.New("EOF Reached while trying to parse whitespace\n")
		}

		for {
			bchar, _, err := ctx.Stream.ReadRune()
			if err == io.EOF {
				return nil
			}

			if (!strings.ContainsRune(reg, bchar)) {
				ctx.Stream.UnreadRune()
				return nil
			}
		}
	}
}

func InitParserContext(parser Parser, root string, stream *strings.Reader) *ParserContext {
	ctx := &ParserContext{
		Stack: make([]string, 0),
		Parser: parser,
		Stream: stream,
		Processed: true,
		Next: false,
		Ast: &Node{Type: "expr", Childs: make([]*Node, 0)},
		Captured: make([]*strings.Builder, 0),
	}

	ctx.Stack = append(ctx.Stack, root)
	ctx.Captured = append(ctx.Captured, &strings.Builder{})
	return ctx
}

func ParseExpr(ctx *ParserContext) error {
	ctx.Processed = true

	for ctx.Processed == true {
		ctx.Next = false
		curr_expr := ctx.Parser[ctx.Stack[len(ctx.Stack)-1]]
		ctx.Processed = false

		curr_rune, _, err := ctx.Stream.ReadRune()
		if err == io.EOF {
			return errors.New(fmt.Sprintf("EOF While parsing\n"))
		}

		for _, entry := range curr_expr {
			if (entry.WildCard || entry.Char == curr_rune) {
				//fmt.Printf("Matched [%s.%s] on rune [%s]\n", ctx.Stack[len(ctx.Stack)-1], entry.Name, string(curr_rune))

				err := entry.Action(ctx, curr_rune)
				if err != nil {
					return err
				}

				if ctx.Next == true {
					break
				}

				if ctx.Stream.Len() == 0 {
					// METTRE LES ERREUR EOF ICI
					ctx.Ast = FindRootNode(ctx.Ast)
					return nil
				}
			}
		}
	}

	return errors.New(fmt.Sprintf("Error parsing all this mess\n"))
}

func GetParser() {
	parser := Parser{
		"expr": {
			Rule{"OpenParent", '(', Push("list"), false},
			Rule{"OpenQuote", '"', Push("string"), false},

			Rule{"comment", ';', Push("comment"), false},
			Rule{"Whitespace", '_', Ignore(" \n\t,"), true},
		},
		"list": {
			Rule{"CloseParent", ')', Pop(), false},
			Rule{"OpenParent", '(', Push("list"), false},
			Rule{"OpenQuote", '"', Push("string"), false},

			Rule{"comment", ';', Push("comment"), false},
			Rule{"Whitespace", '_', Ignore(" \n\t,"), true},
		},
		"comment": {
			Rule{"NewLine", '\n', Pop(), false},
			Rule{"Char", '_', Read(), true},
		},
		"string": {
			Rule{"Escape", '\\', Goto("escaped"), false},
			Rule{"CloseQuote", '"', Pop(), false},
			Rule{"Char", '_', Read(), true},
		},
		"escaped": {
			Rule{"Char", '_', Return(), true},
		},
	}

	data := strings.NewReader(`("lol", ("qsdfqsdfqsd\\  \ \" \\\\" "qsdf" ) "lol") ; this is a comment`)

	ctx := InitParserContext(parser, "expr", data)

	parse_err := ParseExpr(ctx)
	if parse_err != nil {
		fmt.Print(parse_err)
		return
	}

	//fmt.Print("AST:\n")
	//DisplayNode(ctx.Ast, 0)

	bexpr, bexpr_err := ProcessNode(ctx.Ast)
	if bexpr_err != nil {
		fmt.Print(bexpr_err)
		return
	}

	var sb strings.Builder
	disp_err := Display(&bexpr, &sb, true)
	if disp_err != nil {
		fmt.Printf("ERROR %s\n", bexpr_err)
		return
	}

	fmt.Printf("BType expression: %+v\n", sb.String())

	
}

func Unescape(s string) string {
	// strconv.Unquote
	var res strings.Builder
	esc := false

	for _, i := range s {
		if i == '\\' && esc == false {
			esc = true
		} else if i == '\\' && esc == true {
			res.WriteRune('\\')
			esc = false
		} else if i == 'n' && esc == true {
			res.WriteRune('\n')
			esc = false
		} else if i == '"' && esc == true {
			res.WriteRune('"')
			esc = false
		} else {
			res.WriteRune(i)
			esc = false
		}
	}

	return res.String()
}

func Escape(s string) string {
	var res strings.Builder
	for _, i := range s {
		if i == '\\' {
			res.WriteString("\\\\")
		} else if i == '"' {
			res.WriteString("\\\"")
		} else if i == '\n' {
			res.WriteString("\\n")
		} else {
			res.WriteRune(i)
		}
	}

	return res.String()
}

func PrStr(x string, readably bool) string {
     if readably {
	     return Escape(x)
     } else {
	     return x
     }
}

func betterFormat(num float64) string {
    s := fmt.Sprintf("%.6f", num)
    return strings.TrimRight(strings.TrimRight(s, "0"), ".")
}

/*
func Display(x *BType, readably bool) (string, error) {
	if x.BBoolean != nil {
		if x.BBoolean.Value {
			return "true", nil
		} else {
			return "false", nil
		}
	}

	case BFunc:
		return "#<function>"

	if x.BList != nil {
		var res strings.Builder

		for _, s := range x.BList.Values {
			o, err := Display(s, readably)
			if err != nil {
				return "", err
			}
			res.WriteString(o)
			res.WriteString(" ")
		}

		r, size := utf8.DecodeLastRuneInString(res.String())
		if r == utf8.RuneError && (size == 0 || size == 1) {
			size = 0
		}

		return fmt.Sprintf("(%s)", res.String()[:len(res.String())-size]), nil
	}

	if x.BNumber != nil {
		return betterFormat(x.BNumber.Value), nil
	}

	if x.BString != nil {
		if readably {
			return fmt.Sprintf("\"%s\"", PrStr(x.BString.Value, readably)), nil
		} else {
			return x.BString.Value, nil
		}
	}

	if x.BVector != nil {
		var res strings.Builder

		for _, s := range x.BVector.Values {
			o, err := Display(s, readably)
			if err != nil {
				return "", err
			}
			res.WriteString(o)
			res.WriteString(" ")
		}

		r, size := utf8.DecodeLastRuneInString(res.String())
		if r == utf8.RuneError && (size == 0 || size == 1) {
			size = 0
		}

		return fmt.Sprintf("[%s]", res.String()[:len(res.String())-size]), nil
	}

	if x.BHashMap != nil {
		var res strings.Builder

		for _, s := range x.BHashMap.Map {
			o, err := Display(s.Key, readably)
			res.WriteString(o)
			if err != nil {
				return "", err
			}
			res.WriteString(" ")
			o1, err1 := Display(s.Value, readably)
			res.WriteString(o1)
			if err1 != nil {
				return "", err1
			}
			res.WriteString(" ")
		}

		r, size := utf8.DecodeLastRuneInString(res.String())
		if r == utf8.RuneError && (size == 0 || size == 1) {
			size = 0
		}

		return fmt.Sprintf("{%s}", res.String()[:len(res.String())-size]), nil
	}
	
	if x.BKeyword != nil {
		return fmt.Sprintf(":%s", x.BKeyword.Value), nil
	}

	if x.BUnquote != nil {
		o3, err3 := Display(x.BUnquote.Value, readably)
		if err3 != nil {
			return "", err3
		}

		return fmt.Sprintf("'%s", o3), nil
	}

	if x.BQuasiquote != nil {
		o2, err2 := Display(x.BQuasiquote.Value, readably)
		if err2 != nil {
			return "", err2
		}

		return fmt.Sprintf("`%s", o2), nil
	}

	case BException:
		return Display(val.Message, readably)

	if x.BName != nil {
		return x.BName.Value, nil
	}

	case BAtom:
		return "(atom {})".format(Display(val.data, readably))

	if x.BNil != nil {
		return "nil", nil
	}

	panic(fmt.Sprintf("Unable to display: %+v", x))
}
*/

/*
func ParseFile(data string) (*Program) {
	tree := &Program{}
	err := parser.ParseString(data, tree)
	if err != nil {
		return err, nil
	}

	return nil, tree
}
*/

func Prnt(e *BType) (error) {
	var sb strings.Builder

	err := Display(e, &sb, true)
	if err != nil {
		return err
	}
	fmt.Print(sb.String() + "\n")
	return nil
}
