#!/usr/bin/env python
import types
from lark import Lark, Transformer, Token
from basl_types import Name, Keyword, Fn, Atom

rules=r'''
?start: obj |

?obj: list
    | metadata
    | deref
    | hashmap
    | vector
    | keyword
    | quote
    | quasiquote
    | spliceunquote
    | unquote
    | python
    | COMMA
    | TOKEN -> name
    | COMMENT
    | NUM -> number
    | BOOLEAN -> boolean
    | string
    | variadic
    | NIL -> nil

list: "(" obj* ")"
metadata: "^" obj obj
deref: "@" obj
hashmap: "{" ((keyword|string) obj)* "}"
vector: "[" obj* "]"
keyword: ":" TOKEN
quote: "'" obj
quasiquote: "`" obj
unquote: "~" obj
spliceunquote: "~@" obj
python: "\." TOKEN
string: ESCAPED_STRING
variadic: "&"

NIL.5: "nil"
BOOLEAN.5: /true|false/
NUM.5: "-"?NUMBER

COMMENT: /;.*(?=(\n|$))/
COMMA: ","

TOKEN: /[^"^.@~`\[\]:{}&0-9\s,();][^"^@~`\[\]:{}\s();]*/

%import common.ESCAPED_STRING
%import common.NUMBER
%import common.WS
%ignore WS
%ignore COMMENT
%ignore COMMA
'''

l = Lark(rules, parser='lalr', start="start")

class ToAst(Transformer):
    start = lambda _,x: x[0] if len(x) else None
    list = tuple
    vector = lambda _,x: list(x)

    nil = lambda _,x: None
    variadic = lambda _,x: Name("&")
    number = lambda _,x: float(x[0].value) if x[0].value.find(".") != -1 else int(x[0].value) 
    boolean = lambda _,x: x[0] == "true"
    name = lambda _,x: Name(x[0].value)
    string = lambda _,x: eval(x[0])
    deref = lambda _,x: tuple([Name("deref"), *x])
    metadata = lambda _,x: tuple([Name("with-meta"), x[0], x[1]])
    hashmap = lambda _,x: { i[0]: i[1] for i in zip(list(x[::2]), list(x[1::2])) }
    keyword = lambda _,x: Keyword(x[0].value)
    quote = lambda _,x: tuple([Name("quote"), *x])
    quasiquote = lambda _,x: tuple([Name("quasiquote"), *x])
    unquote = lambda _,x: tuple([Name("unquote"), *x])
    spliceunquote = lambda _,x: tuple([Name("spliceunquote"), *x])

def display(x, print_readably=True):
    if isinstance(x, bool):
        return "true" if x is True else "false"

    if isinstance(x, types.LambdaType):
        return "#<function>"

    if isinstance(x, Fn):
        return "#<function>"

    if isinstance(x, tuple):
        return "({})".format(" ".join([display(r, print_readably) for r in x]))

    if isinstance(x, int):
        return repr(x)

    if isinstance(x, float):
        return repr(x)

    if isinstance(x, str):
        return repr(x) if print_readably else x

    if isinstance(x, list):
        return "[{}]".format(" ".join([display(s, print_readably) for s in x]))

    if isinstance(x, dict):
        return "{{{}}}".format(" ".join(["{} {}".format(":{}".format(k), display(v, print_readably)) for k,v in x.items()]))

    if isinstance(x, Keyword):
        return ":{}".format(x.name)

    if isinstance(x, Name):
        return x.name

    if isinstance(x, Atom):
        return "(atom {})".format(display(x.data))

    if x is None:
        return "nil"

    return x

def parse(data):
    tree = l.parse(data)
    return ToAst().transform(tree)

if __name__ == "__main__":
    if len(sys.argv) >= 2:
        [print(display(a)) for a in parse(open(sys.argv[1], "r").read())]

