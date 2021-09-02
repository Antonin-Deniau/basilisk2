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
	Text string
	Parser Parser
	Index int64
	Ast *Node
}


// TO IMPLEMENT
// metadata: "^" obj obj
// deref: "@" obj
// hashmap: "{" ((keyword|string) obj)* "}"
// vector: "[" obj* "]"
// quasiquote: "`" obj
// unquote: "~" obj
// spliceunquote: "~@" obj
// variadic: "&"

var string_regex = regexp.MustCompile(`^"((:?\\.|[^"\\])*)"`)
var name_regex = regexp.MustCompile("^([^\"^.@~`\\[\\]:{}'0-9\\s,();][^\"^@~`\\[\\]:{}\\s();]*)")
var open_parent_regex = regexp.MustCompile(`^\(`)
var close_parent_regex = regexp.MustCompile(`^\)`)
var keyword_regex = regexp.MustCompile("^:([^\"^.@~`\\[\\]:{}'0-9\\s,();][^\"^@~`\\[\\]:{}\\s();]*)")
var whitespace_regex = regexp.MustCompile(`^((;[^\n]*$)|[\s,]+)`)
var bool_regex = regexp.MustCompile(`^(true|false)\b`)
var int_regex = regexp.MustCompile(`^(-?[0-9]+)`)
var nil_regex = regexp.MustCompile(`^nil\b`)
var quote_regex = regexp.MustCompile(`^'`)

var rules = Parser{
	"Expr": {
		Rule{"List", open_parent_regex, Push("List")},
		Rule{"Quote", quote_regex, PushOne("Expr")},

		Rule{"Nil", nil_regex, ReadRegex()},
		Rule{"Bool", bool_regex, ReadRegex()},
		Rule{"String", string_regex, ReadRegex()},
		Rule{"Int", int_regex, ReadRegex()},
		Rule{"Name", name_regex, ReadRegex()},
		Rule{"Keyword", keyword_regex, ReadRegex()},
	},
	"List": {
		Rule{"List", open_parent_regex, Push("List")},
		Rule{"Quote", quote_regex, PushOne("Expr")},

		Rule{"Nil", nil_regex, ReadRegex()},
		Rule{"Bool", bool_regex, ReadRegex()},
		Rule{"String", string_regex, ReadRegex()},
		Rule{"Int", int_regex, ReadRegex()},
		Rule{"Name", name_regex, ReadRegex()},
		Rule{"Keyword", keyword_regex, ReadRegex()},
		Rule{"EndList", close_parent_regex, Pop()},
	},
}

func Ignore(ctx *ParserContext) {
	found := whitespace_regex.FindStringSubmatch(ctx.Text[ctx.Index:])
	if found != nil {
		ctx.Index += int64(len(found[0]))
	}
}

func Push(next_expr string) Action {
	return func(ctx *ParserContext, matcher_name string, captured string) error {
		child := &Node{
			Type: matcher_name,
			Childs: make([]*Node, 0),
			Repeat: true,
			Validated: false,
			ParserRule: next_expr,
			Parent: ctx.Ast,
		}

		ctx.Ast.Childs = append(ctx.Ast.Childs, child)
		ctx.Ast = child

		return nil
	}
}


func PushOne(next_expr string) Action {
	return func(ctx *ParserContext, matcher_name string, captured string) error {
		child := &Node{
			Type: matcher_name,
			Repeat: false,
			Validated: false,
			ParserRule: next_expr,
			Childs: make([]*Node, 0),
			Parent: ctx.Ast,
		}

		ctx.Ast.Childs = append(ctx.Ast.Childs, child)
		ctx.Ast = child

		return nil
	}
}

func ReadRegex() Action {
	return func(ctx *ParserContext, matcher_name string, captured string) error {
		child := &Node{
			Type: matcher_name,
			Parent: ctx.Ast,
			Value: captured,
			Validated: true,
		}

		ctx.Ast.Childs = append(ctx.Ast.Childs, child)

		return nil
	}
}

func Pop() Action {
	return func(ctx *ParserContext, matcher_name string, captured string) error {
		ctx.Ast.Validated = true
		ctx.Ast = ctx.Ast.Parent

		return nil
	}
}

func InitParserContext(str string) *ParserContext {
	ctx := &ParserContext{
		Parser: rules,
		Index: 0,
		Text: str,
		Ast: &Node{Type: "Expr", Validated: false, Repeat: false, ParserRule: "Expr", Childs: make([]*Node, 0)},
	}

	return ctx
}

func ParseExpr(ctx *ParserContext) error {
	matched := false

	for {
		matched = false

		curr_expr := ctx.Parser[ctx.Ast.ParserRule]
		Ignore(ctx)

		for _, entry := range curr_expr {
	    	found := entry.Regex.FindStringSubmatch(ctx.Text[ctx.Index:])
			//fmt.Printf("Check [%s.%s] == %s\n", ctx.Ast.Type, entry.Name, found)

			if found != nil {
				matched = true

				ctx.Index += int64(len(found[0]))

				var test_str string
				if len(found) != 1 {
					test_str = found[1]
				} else {
					test_str = ""
				}

				err_action := entry.Action(ctx, entry.Name, test_str)
				if err_action != nil {
					return err_action
				}

				//fmt.Printf("matched => %+v\n", found[0])
				//fmt.Printf("======================\n")

				break
			}
		}


		// >>>IGNORE EOF WHITESPACES
		if ctx.Ast.Parent == nil {
			Ignore(ctx)
		}

		if matched == false {
			if int64(len(ctx.Text)) != ctx.Index  {
				break
			} else {
				return nil
			}
		}
		// <<<IGNORE EOF WHITESPACES

		if ctx.Ast.Repeat == true {
			continue
		} else {
			if ctx.Ast.Validated == true {
				if ctx.Ast.Parent == nil {

					return nil
				} else {
					ctx.Ast = ctx.Ast.Parent
					continue
				}
			}

			ctx.Ast.Validated = true

			continue
		}
	}

	return errors.New(fmt.Sprintf("Error cannot parse, stopped at =>%s\n", ctx.Text[ctx.Index:]))
}

func Parse(str_input string) (*BType, error) {
	ctx := InitParserContext(str_input)

	parse_err := ParseExpr(ctx)
	DisplayNode(ctx.Ast, 2)

	if parse_err != nil {
		return nil, parse_err
	}

	bexpr, bexpr_err := ProcessNode(ctx.Ast)
	if bexpr_err != nil {
		return nil, bexpr_err
	}

	return &bexpr, nil
}
