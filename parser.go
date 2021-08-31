package main

import (
	"strings"
	"errors"
	"regexp"
	"fmt"
)

type Rule struct {
	Name string
	Regex *regexp.Regexp
	Action Action
}

type Action func(ctx *ParserContext, matcher_name string, captured string) error

type Parser map[string][]Rule

type ParserContext struct {
	Stack []string
	Text string
	Parser Parser
	Processed bool
	Index int64
	Ast *Node
}


var string_regex = regexp.MustCompile(`^"((:?\\.|[^"\\])*)"`)
var name_regex = regexp.MustCompile("^([^\"^.@~`\\[\\]:{}'0-9\\s,();][^\"^@~`\\[\\]:{}\\s();]*)")
var open_parent_regex = regexp.MustCompile(`^\(`)
var close_parent_regex = regexp.MustCompile(`^\)`)
var keyword_regex = regexp.MustCompile("^:([^\"^.@~`\\[\\]:{}'0-9\\s,();][^\"^@~`\\[\\]:{}\\s();]*)")
var comment_regex = regexp.MustCompile(`^;[^\n]*\n`)
var whitespace_regex = regexp.MustCompile(`^[\s,]+`)


var rules = Parser{
	"Expr": {
		Rule{"OpenParent", open_parent_regex, Push("List")},

		Rule{"String", string_regex, Read()},
		Rule{"Keyword", keyword_regex, Read()},
		Rule{"Name", name_regex, Read()},

		Rule{"Comment", comment_regex, Read()},
		Rule{"Whitespace", whitespace_regex, Read()},
	},
	"List": {
		Rule{"CloseParent", close_parent_regex, Pop()},
		Rule{"OpenParent", open_parent_regex, Push("List")},

		Rule{"String", string_regex, Read()},
		Rule{"Name", name_regex, Read()},
		Rule{"Keyword", keyword_regex, Read()},

		Rule{"Comment", comment_regex, Read()},
		Rule{"Whitespace", whitespace_regex, Read()},
	},
}



func Push(next_expr string) Action {
	return func(ctx *ParserContext, matcher_name string, captured string) error {
		//fmt.Printf("PUSH: %s\n", next_expr)

		ctx.Stack = append(ctx.Stack, next_expr)

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
	return func(ctx *ParserContext, matcher_name string, captured string) error {
		//fmt.Printf("READ: [%s]\n", captured)

		child := &Node{
			Type: matcher_name,
			Parent: ctx.Ast,
			Value: captured,
		}

		ctx.Ast.Childs = append(ctx.Ast.Childs, child)

		return nil
	}
}

func Pop() Action {
	return func(ctx *ParserContext, matcher_name string, captured string) error {
		//fmt.Printf("POP [%s]\n", matcher_name)

		ctx.Stack = ctx.Stack[:len(ctx.Stack)-1]

		ctx.Ast = ctx.Ast.Parent

		return nil
	}
}

func InitParserContext(str string) *ParserContext {
	ctx := &ParserContext{
		Stack: make([]string, 0),
		Parser: rules,
		Processed: true,
		Index: 0,
		Text: str,
		Ast: &Node{Type: "Expr", Childs: make([]*Node, 0)},
	}

	ctx.Stack = append(ctx.Stack, "Expr")
	return ctx
}

func ParseExpr(ctx *ParserContext) error {
	ctx.Processed = true

	for {
		curr_expr := ctx.Parser[ctx.Stack[len(ctx.Stack)-1]]
		ctx.Processed = false

		for _, entry := range curr_expr {
	    	found := entry.Regex.FindStringSubmatch(ctx.Text[ctx.Index:])

	    	//fmt.Printf("Checking [%s.%s] => %s\n", ctx.Stack[len(ctx.Stack)-1], entry.Name, ctx.Text[ctx.Index:])

			if found != nil {
				ctx.Index += int64(len(found[0]))

				if len(found) != 1 {
					//fmt.Printf("Matched [%s.%s] [%s]\n", ctx.Stack[len(ctx.Stack)-1], entry.Name, found[1])

					err_action := entry.Action(ctx, entry.Name, found[1])
					if err_action != nil {
						return err_action
					}
				} else {
					//fmt.Printf("Matched [%s.%s] _\n", ctx.Stack[len(ctx.Stack)-1], entry.Name)
					err_action := entry.Action(ctx, entry.Name, "")
					if err_action != nil {
						return err_action
					}
				}

				if int64(len(ctx.Text)) == ctx.Index {
					if len(ctx.Stack) != 1 {
						return errors.New("Unexpected EOF")
					}

					ctx.Ast = FindRootNode(ctx.Ast)
					return nil
				}

				break
			}
		}
	}

	return errors.New(fmt.Sprintf("Error parsing all this mess\n"))
}

func TestParser(str_input string) error {
	ctx := InitParserContext(fmt.Sprintf("%s\n", str_input))

	parse_err := ParseExpr(ctx)
	if parse_err != nil {
		return parse_err
	}

	DisplayNode(ctx.Ast, 0)

	return nil
}


func Parse(str_input string) (*BType, error) {
	ctx := InitParserContext(fmt.Sprintf("%s\n", str_input))

	parse_err := ParseExpr(ctx)
	if parse_err != nil {
		return nil, parse_err
	}

	bexpr, bexpr_err := ProcessNode(ctx.Ast)
	if bexpr_err != nil {
		return nil, bexpr_err
	}

	return &bexpr, nil
	/*
	var sb strings.Builder
	disp_err := Display(&bexpr, &sb, true)
	if disp_err != nil {
		fmt.Printf("ERROR %s\n", bexpr_err)
		return
	}
	fmt.Printf("BType expression: %+v\n", sb.String())
	*/	
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
