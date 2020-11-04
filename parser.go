package main

import (
	"github.com/alecthomas/participle/lexer"
)

type Program struct {
	Pos lexer.Position

	Lines []*BType `@@*`
}

graphQLLexer = lexer.Must(stateful.New(stateful.Rules{
	"Root": {
		{"Comment", `;.*(?=(\n|$))`, nil },
		{"Ident", "[^\"^.@~`\\[\\]:{}&'0-9\\s,();][^\"^@~`\\[\\]:{}\\s();]*", nil },
		{"Number", `-?\d+(\.\d+)?`, nil },
		{"Whitespace", `(\s|\t|\n|,|\r)+`, nil },
		{"String", `"(\.|[^"\\])*"`, nil },
		{"Nil", `nil`, nil },
		{"Bool", `true|false`, nil },
	}
))

parser = participle.MustBuild(&Program{},
    participle.Lexer(graphQLLexer),
    participle.Elide("Comment", "Whitespace"),
)

func Unescape(s string) string {
	res :=  ""
	esc := false

	for _, i in range s {
		if i == '\\' && esc == false {
			esc = true
		} else if i == '\\' && esc == true {
			res += "\\"
			esc = false
		} else if i == "n" && esc == true {
			res += "\n"
			esc = false
		} else if i == '"' && esc == true {
			res += '"'
			esc = false
		} else {
			res += i
			esc = false
		}
	}

	return res
}

func Escape(s string) string {
	res =  ""
	for _, i in range s {
		if i == "\\" {
			res += "\\\\"
		} else if i == '"' {
			res += '\\"'
		} else if i == '\n' {
			res += '\\n'
		} else {
			res += i
		}
	}

	return res
}

func PrStr(x string, readably bool) string {
     if readably {
	     return Escape(x)
     } else {
	     return x
     }
}

func Display(x *BType, readably bool) string {
    if isinstance(x, bool):
        return "true" if x is True else "false"

    if isinstance(x, types.LambdaType):
        return "#<function>"

    if isinstance(x, Fn):
        return "#<function>"

    if isinstance(x, tuple):
        return "({})".format(" ".join([display(r, readably) for r in x]))

    if isinstance(x, int):
        return repr(x)

    if isinstance(x, float):
        return repr(x)

    if isinstance(x, str):
        return "\"{}\"".format(pr_str(x, readably)) if readably else x

    if isinstance(x, list):
        return "[{}]".format(" ".join([display(s, readably) for s in x]))

    if isinstance(x, dict):
        return "{{{}}}".format(
                " ".join(["{} {}".format(display(k, readably), display(v, readably)) for k,v in x.items()]))

    if isinstance(x, Keyword):
        return ":{}".format(x.name)

    if isinstance(x, BaslException):
        return display(x.message, readably)

    if isinstance(x, Name):
        return x.name

    if isinstance(x, Atom):
        return "(atom {})".format(display(x.data, readably))

    if x is None:
        return "nil"

    return x
}

func Parse(data) {
    tree = l.parse(data)
    return ToAst().transform(tree)
}

func Prnt(e) {
    sys.stdout.write(display(e, True))
    sys.stdout.write("\n")
}

func Input(repl string) {
}

func main() {
    for {
	res, err = Prnt(Parse(Input("basilisk> ")))
	if err != nil {
		print("EOF: {}".format(e))
		continue
	}

	print(res if res != nil else "nil")
    }
}
