#!/usr/bin/env python
from lark import Lark, Transformer, Token

rules=r'''
lines: obj*

?obj: (list|atom|metadata|deref|hashmap|vector|keyword|quote|quasiquote|spliceunquote|unquote|python|COMMA)

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

?atom: name
     | COMMENT
     | NUMBER -> number
     | string
     | "nil" -> nil
     | BOOLEAN -> boolean

BOOLEAN: /true|false/

string: ESCAPED_STRING
name: TOKEN
COMMENT: /;.*\n/
COMMA: ","

TOKEN: /[^"^.@~`\[\]:{}0-9\s,();][^"^@~`\[\]:{}0-9\s();]*/

%import common.ESCAPED_STRING
%import common.NUMBER
%import common.WS
%ignore WS
%ignore COMMENT
%ignore COMMA
'''

l = Lark(rules, parser='lalr', start="lines")

class Name:
    def __init__(self, name):
        self.name = name

    def __hash__(self):
        return hash(self.name)

    def __repr__(self):
        return self.name

    def __eq__(self, a, b):
        return a.name == b.name

class Keyword:
    def __init__(self, name):
        self.name = name

    def __hash__(self):
        return hash(self.name)

    def __repr__(self):
        return self.name

    def __eq__(self, a, b):
        return a.name == b.name

class ToAst(Transformer):
    lines = tuple
    list = tuple

    nil = lambda _,x: None
    number = lambda _,x: float(x[0].value) if x[0].value.find(".") != -1 else int(x[0].value) 
    boolean = lambda _,x: x[0] == "true"
    name = lambda _,x: Name(x[0].value)
    string = lambda _,x: eval(x[0])
    deref = lambda _,x: tuple(Name("deref"), *x )
    metadata = lambda _,x: tuple(Name("with-meta"), x[0], x[1])
    hashmap = lambda _,x: { i[0]: i[1] for i in zip(list(x[::2]), list(x[1::2])) }
    keyword = lambda _,x: Keyword(x[0].value)
    vector = lambda _,x: x
    quote = lambda _,x: tuple(Name("quote"), *x)
    quasiquote = lambda _,x: tuple(Name("quasiquote"), *x)
    unquote = lambda _,x: tuple(Name("unquote"), *x)
    spliceunquote = lambda _,x: tuple(Name("spliceunquote"), *x)

def display(x):
    if isinstance(x, tuple):
        return "({})".format([display(r) for r in x])

    if isinstance(x, int):
        return repr(x)

    if isinstance(x, float):
        return repr(x)

    if isinstance(x, str):
        return repr(x)

    if isinstance(x, list):
        return "[{}]".format(" ".join([display(s) for s in x]))

    if isinstance(x, dict):
        return "{{{}}}".format(" ".join(["{} {}".format(":{}".format(k), display(v)) for k,v in x.items()]))

    if isinstance(x, Keyword):
        return ":{}".format(x.name)

    if isinstance(x, Name):
        return x

    return x

def parse(data):
    tree = l.parse(data)
    return ToAst().transform(tree)

if __name__ == "__main__":
    [print(a) for a in parse(open("./syntax.cr", "r").read())]

