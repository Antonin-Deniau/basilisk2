#!/usr/bin/env python
import types, json, sys, re
from lark import Lark, Transformer, Token
from lark import UnexpectedInput, UnexpectedToken
from basl_types import Name, Keyword, Fn, Atom, BaslException

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

NIL.5: /nil(?!\?)/
BOOLEAN.5: /(true|false)(?!\?)/
NUM.5: "-"?NUMBER

COMMENT: /;.*(?=(\n|$))/
COMMA: ","

TOKEN: /[^"^.@~`\[\]:{}'0-9\s,();][^"^@~`\[\]:{}\s();]*/

ESCAPED_STRING: /"(\\.|[^"\\])*"/

%import common.NUMBER
%import common.WS
%ignore WS
%ignore COMMENT
%ignore COMMA
'''

l = Lark(rules, parser='lalr', start="start")

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

class ToAst(Transformer):
    start = lambda _,x: x[0] if len(x) else None
    list = tuple
    #vector = lambda _,x: list(x)
    vector = lambda _,x: tuple([Name("vector"), *x])

    nil = lambda _,x: None
    variadic = lambda _,x: Name("&")
    number = lambda _,x: float(x[0].value) if x[0].value.find(".") != -1 else int(x[0].value) 
    boolean = lambda _,x: x[0] == "true"
    name = lambda _,x: Name(x[0].value)
    string = lambda _, x: unescape(x[0][1:-1])
    deref = lambda _,x: tuple([Name("deref"), *x])
    metadata = lambda _,x: tuple([Name("with-meta"), x[1], x[0]])
    #hashmap = lambda _,x: { i[0]: i[1] for i in zip(list(x[::2]), list(x[1::2])) }
    hashmap = lambda _,x: tuple([Name("hash-map"), *x])
    keyword = lambda _,x: Keyword(x[0].value)
    quote = lambda _,x: tuple([Name("quote"), *x])
    quasiquote = lambda _,x: tuple([Name("quasiquote"), *x])
    unquote = lambda _,x: tuple([Name("unquote"), *x])
    spliceunquote = lambda _,x: tuple([Name("splice-unquote"), *x])

def parse(data):
    tree = l.parse(data)
    a =  ToAst().transform(tree)
    return a


symbol_dict = []
name_dict = []

def search_symbol(ast):
    global symbol_dict

    if isinstance(ast, Keyword) and ast not in symbol_dict:
        symbol_dict.append(str(ast))
        return
    
    if isinstance(ast, tuple):
        [search_symbol(d) for d in ast]


def search_name(ast):
    global name_dict

    if isinstance(ast, Name) and ast not in name_dict:
        name_dict.append(str(ast))
        return
    
    if isinstance(ast, tuple):
        [search_name(d) for d in ast]


def serialise_list():
    pass

def serialise_number():
    pass

data = parse("(do {})".format(open(sys.argv[1]).read()))

search_name(data)
search_symbol(data)


print(name_dict)
print(symbol_dict)


















