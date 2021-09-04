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

// types
// hashmap: "{" ((keyword|string) obj)* "}"
// vector: "[" obj* "]"

// sugars
// deref: "@" obj
// quasiquote: "`" obj
// unquote: "~" obj
// spliceunquote: "~@" obj

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
var variadic_regex = regexp.MustCompile(`^&`)
var meta_regex = regexp.MustCompile(`^\^`)
var hashmap_open_regex = regexp.MustCompile(`^{`)
var hashmap_close_regex = regexp.MustCompile(`^}`)

var rules = Parser{
	"Expr": {
		Rule{"List", open_parent_regex, Push("List", -1)},
		Rule{"Quote", quote_regex, Push("Expr", 1)},
		Rule{"Meta", meta_regex, Push("Expr", 2)},
		Rule{"Hashmap", hashmap_open_regex, Push("Hashmap", -1)},

		Rule{"Nil", nil_regex, ReadRegex()},
		Rule{"Bool", bool_regex, ReadRegex()},
		Rule{"String", string_regex, ReadRegex()},
		Rule{"Int", int_regex, ReadRegex()},
		Rule{"Name", name_regex, ReadRegex()},
		Rule{"Keyword", keyword_regex, ReadRegex()},
		Rule{"Variadic", variadic_regex, ReadRegex()},
	},
	"Hashmap": {
		Rule{"List", open_parent_regex, Push("List", -1)},
		Rule{"Quote", quote_regex, Push("Expr", 1)},
		Rule{"Meta", meta_regex, Push("Expr", 2)},
		Rule{"Hashmap", hashmap_open_regex, Push("Hashmap", -1)},

		Rule{"Nil", nil_regex, ReadRegex()},
		Rule{"Bool", bool_regex, ReadRegex()},
		Rule{"String", string_regex, ReadRegex()},
		Rule{"Int", int_regex, ReadRegex()},
		Rule{"Name", name_regex, ReadRegex()},
		Rule{"Keyword", keyword_regex, ReadRegex()},
		Rule{"Variadic", variadic_regex, ReadRegex()},

		Rule{"EndMap", hashmap_close_regex, Pop()},
	},
	"List": {
		Rule{"List", open_parent_regex, Push("List", -1)},
		Rule{"Quote", quote_regex, Push("Expr", 1)},
		Rule{"Meta", meta_regex, Push("Expr", 2)},
		Rule{"Hashmap", hashmap_open_regex, Push("Hashmap", -1)},

		Rule{"Nil", nil_regex, ReadRegex()},
		Rule{"Bool", bool_regex, ReadRegex()},
		Rule{"String", string_regex, ReadRegex()},
		Rule{"Int", int_regex, ReadRegex()},
		Rule{"Name", name_regex, ReadRegex()},
		Rule{"Keyword", keyword_regex, ReadRegex()},
		Rule{"Variadic", variadic_regex, ReadRegex()},

		Rule{"EndList", close_parent_regex, Pop()},
	},
}

func Ignore(ctx *ParserContext) {
	found := whitespace_regex.FindStringSubmatch(ctx.Text[ctx.Index:])
	if found != nil {
		ctx.Index += int64(len(found[0]))
	}
}

func Push(next_expr string, limit int) Action {
	return func(ctx *ParserContext, matcher_name string, captured string) error {
		child := &Node{
			Type: matcher_name,
			Childs: make([]*Node, 0),
			Limit: limit,
			Counter: 0,
			ParserRule: next_expr,
			Parent: ctx.Ast,
		}
		
		ctx.Ast.Counter += 1

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
			Limit: 1,
			Counter: 1,
		}

		ctx.Ast.Counter += 1
		ctx.Ast.Childs = append(ctx.Ast.Childs, child)

		return nil
	}
}

func Pop() Action {
	return func(ctx *ParserContext, matcher_name string, captured string) error {
		ctx.Ast = ctx.Ast.Parent

		return nil
	}
}

func InitParserContext(str string) *ParserContext {
	ctx := &ParserContext{
		Parser: rules,
		Index: 0,
		Text: str,
		Ast: &Node{
			Type: "Expr",
			ParserRule: "Expr",
			Limit: 1,
			Counter: 0,
			Childs: make([]*Node, 0),
		},
	}

	return ctx
}

func ParseExpr(ctx *ParserContext) error {
	matched := false

	for {
		Ignore(ctx)

		if ctx.Ast.Limit != -1 && ctx.Ast.Counter >= ctx.Ast.Limit {
			//fmt.Println("POP")
			//fmt.Printf("%+v\n", ctx.Ast)
			if ctx.Ast.Parent == nil {
				return nil
			}

			ctx.Ast = ctx.Ast.Parent
		}

		curr_expr := ctx.Parser[ctx.Ast.ParserRule]

		matched = false
		for _, entry := range curr_expr {
	    	found := entry.Regex.FindStringSubmatch(ctx.Text[ctx.Index:])
			//fmt.Printf("Check [%s.%s] == %s\n", ctx.Ast.Type, entry.Name, found)
			//fmt.Printf("Check %s\n", ctx.Text[ctx.Index:])
			//fmt.Printf("%n\n", ctx.Ast.Counter)
			//fmt.Printf("%n\n", ctx.Ast.Limit)

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

				//fmt.Printf("matched => [%s]\n", found[0])
				//fmt.Printf("Ns [%s.%s]\n", ctx.Ast.Type, entry.Name)
				//fmt.Printf("Quantity %d/%d\n", ctx.Ast.Counter, ctx.Ast.Limit)
	    		//fmt.Println("========================================")
				break
			}
		}

		if matched == false {
			break
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
