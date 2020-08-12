#!/usr/bin/env python
import json, re, sys
from lark import Lark, Transformer, Token

rawdata = open(sys.argv[1], "r").read()

rules=r'''
lines: obj*

obj: (list|atom|metadata|deref|hashmap|vector|keyword|quote|quasiquote|spliceunquote|unquote|python|COMMA)

list: "(" obj* ")"
metadata: "^" obj obj
deref: "@" obj
hashmap: "{" ((string|keyword|quote) obj)* "}"
vector: "[" obj* "]"
keyword: ":" name
quote: "'" obj
quasiquote: "`" obj
unquote: "~" obj
spliceunquote: "~@" obj
python: "\." name

atom: name
     | COMMENT
     | NUMBER -> number
     | string
     | "nil" -> is_nil
     | "null" -> is_null
     | "false" -> is_false
     | "true" -> is_true

string: ESCAPED_STRING
name: /[^"^.@~`\[\]:{}0-9\s();]+/
COMMENT: /;.*\n/
COMMA: ","

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
    expr = list

    name = lambda _,x: { "type": "name", "value": x[0].value }
    string = lambda _,x: { "type": "string", "value": eval(x[0].value) }
    number = lambda _,x: {
            "type": "number",
            "value":  float(x[0].value) if x[0].value.find(".") else int(x[0].value),
    }

tree = l.parse(rawdata)
lines = ToAst().transform(tree)

for line in lines:
    print(line)
exit()
