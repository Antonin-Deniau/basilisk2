#!/usr/bin/env python
from lark import Lark, Transformer, Token

rules=r'''
lines: obj*

?obj: (list|atom|metadata|deref|hashmap|vector|keyword|quote|quasiquote|spliceunquote|unquote|python|COMMA)

list: "(" obj* ")"
metadata: "^" obj obj
deref: "@" obj
hashmap: "{" ((string|keyword|quote) obj)* "}"
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

class ToAst(Transformer):
    lines = list
    list = list


    nil = lambda _,x: { "type": "nil", "value": None  }
    boolean = lambda _,x: { "type": "boolean", "value": True if x[0].value == "true" else False  }
    name = lambda _,x: { "type": "name", "value": x[0].value }
    string = lambda _,x: { "type": "string", "value": eval(x[0].value) }
    number = lambda _,x: {
            "type": "number",
            "value":  float(x[0].value) if x[0].value.find(".") != -1 else int(x[0].value),
    }
    deref = lambda _,x: [{"type": "name", "value": "deref" }, *x ]
    metadata = lambda _,x: [{"type": "name", "value": "with-meta" }, x[0], x[1] ]
    hashmap = lambda _,x: {
            "type": "hashmap",
            "value": [ [x[0], x[1]] for i in zip(list(x[::2]), list(x[1::2])) ],
    }
    keyword = lambda _,x: {"type": "keyword", "value": x[0].value }
    vector = lambda _,x: {"type": "vector", "value": x }
    quote = lambda _,x: [{"type": "name", "value": "quote"}, *x]
    quasiquote = lambda _,x: [{"type": "name", "value": "quasiquote"}, *x]
    unquote = lambda _,x: [{"type": "name", "value": "unquote"}, *x]
    spliceunquote = lambda _,x: [{"type": "name", "value": "spliceunquote"}, *x]

display_funcs = {
    "nil": lambda x: "nil",
    "boolean": lambda x: "true" if x is True else "false",
    "name": lambda x: x,
    "string": lambda x: repr(x),
    "number": lambda x: repr(x),
    "hashmap": lambda x: "{{{}}}".format(" ".join(["{}:{}".format(display(i[0]), display(i[1])) for i in x])),
    "keyword": lambda x: ":{}".format(x),
    "vector": lambda x: "[%s]" % " ".join([display(i) for i in x]),
}

def display(x):
    if isinstance(x, list):
        return "({})".format(" ".join([display(i) for i in x]))
    else:
        return display_funcs[x["type"]](x["value"])

def parse(data):
    tree = l.parse(data)
    return ToAst().transform(tree)

