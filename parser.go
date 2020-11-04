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

def unescape(s):
    res =  ""
    esc = False
    for i in s:
        if i == '\\' and esc == False:
            esc = True
        elif i == '\\' and esc == True:
            res += "\\"
            esc = False
        elif i == "n" and esc == True:
            res += "\n"
            esc = False
        elif i == '"' and esc == True:
            res += '"'
            esc = False
        else:
            res += i
            esc = False

    return res

def escape(s):
    res =  ""
    for i in s:
        if i == "\\":
            res += "\\\\"
        elif i == '"':
            res += '\\"'
        elif i == '\n':
            res += '\\n'
        else:
            res += i

    return res

def pr_str(x, readably):
    return escape(x) if readably else x

def display(x, readably):
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

def parse(data):
    tree = l.parse(data)
    return ToAst().transform(tree)

def prnt(e):
    sys.stdout.write(display(e, True))
    sys.stdout.write("\n")


if __name__ == "__main__":
    while True:
        try:
            res = prnt(parse(input("basilisk> ")))
            print(res if res != None else "nil")
        except Exception as e:
            print("EOF: {}".format(e))
