package main

import (
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
var bool_regex = regexp.MustCompile(`^(true|false)\b`)
var nil_regex = regexp.MustCompile(`^nil\b`)
var quote_regex = regexp.MustCompile(`^'`)
var nop_regex = regexp.MustCompile(`^`)


var rules = Parser{
	"Expr": {
		Rule{"Comment", comment_regex, Read()},
		Rule{"Whitespace", whitespace_regex, Read()},

		Rule{"List", open_parent_regex, Push("List")},
		Rule{"Quote", quote_regex, Push("Expr")},

		Rule{"Nil", nil_regex, Read()},
		Rule{"Bool", bool_regex, Read()},
		Rule{"String", string_regex, Read()},
		Rule{"Keyword", keyword_regex, Read()},
		Rule{"Name", name_regex, Read()},
	},
	"List": {
		Rule{"EndList", close_parent_regex, Pop()},
		Rule{"Expr", nop_regex, Push("Expr")},
	},
}



func Push(next_expr string) Action {
	return func(ctx *ParserContext, matcher_name string, captured string) error {
		//fmt.Printf("PUSH: %s\n", next_expr)

		ctx.Stack = append(ctx.Stack, next_expr)

		child := &Node{
			Type: matcher_name,
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

		ctx.Stack = ctx.Stack[:len(ctx.Stack)-1]

		if ctx.Ast.Parent != nil {
			ctx.Ast = ctx.Ast.Parent
		}

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
		Index: 0,
		Text: str,
		Ast: &Node{Type: "Expr", Childs: make([]*Node, 0)},
	}

	ctx.Stack = append(ctx.Stack, "Expr")
	return ctx
}

func ParseExpr(ctx *ParserContext) error {
	processed := false
	for {
		processed = false
		curr_expr := ctx.Parser[ctx.Stack[len(ctx.Stack)-1]]

		for _, entry := range curr_expr {
	    	found := entry.Regex.FindStringSubmatch(ctx.Text[ctx.Index:])

	    	//fmt.Printf("Checking [%s.%s] => %s\n", ctx.Stack[len(ctx.Stack)-1], entry.Name, ctx.Text[ctx.Index:])

			if found != nil {
				//fmt.Printf("======================\n")
				//fmt.Printf("matched => %+v\n", found)
				//fmt.Printf("m_len   => %+v\n", int64(len(found[0])))
				//fmt.Printf("index   => %+v\n", ctx.Index)
				//fmt.Printf("Matched [%s.%s] [%s]\n", ctx.Stack[len(ctx.Stack)-1], entry.Name, found[0])

				processed = true
				ctx.Index += int64(len(found[0]))

				if len(found) != 1 {
					err_action := entry.Action(ctx, entry.Name, found[1])
					if err_action != nil {
						return err_action
					}
				} else {
					err_action := entry.Action(ctx, entry.Name, "")
					if err_action != nil {
						return err_action
					}
				}

				if len(ctx.Stack) == 0 {
					if len(ctx.Stack) != 0 {
						return errors.New("Unexpected EOF")
					}

					ctx.Ast = FindRootNode(ctx.Ast)
					return nil
				}

				break
			}
		}

		if processed == false {
			break
		}
	}

	return errors.New(fmt.Sprintf("Error parsing all this mess\n"))
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
}
