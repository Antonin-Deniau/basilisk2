from parser import display, parse
from basl_types import Name

def prn(*a):
    print(" ".join([display(i) for i in a]))
    return None

def println(*a):
    print(" ".join([display(i, False) for i in a]))
    return None

def equality(a, b):
    if (type(a) == tuple or type(a) == list) and (type(b) == tuple or type(b) == list):
        if len(a) != len(b): return False
        for i in zip(a, b):
            if not equality(i[0], i[1]): return False
        return True

    return type(a) == type(b) and a == b

def read_string(a):
    try:
        return parse(a)
    except IndexError:
        return None

ns = {
    '+': lambda a,b: a+b,
    '-': lambda a,b: a-b,
    '*': lambda a,b: a*b,
    '/': lambda a,b: int(a/b),
    'list': lambda *a: tuple(a),
    'list?': lambda a: isinstance(a,tuple),
    'empty?': lambda a: len(a) == 0,
    'count': lambda a: 0 if a == None else len(a),
    '=': equality,
    '<': lambda a, b: a < b,
    '<=': lambda a, b: a <= b,
    '>=': lambda a, b: a >= b,
    '>': lambda a, b: a > b,
    'pr-str': lambda *a: " ".join([display(i) for i in a]),
    'str': lambda *a: "".join([display(i, False) for i in a]),
    'prn': prn,
    'println': println,
    'read-string': read_string,
    'slurp': lambda a: open(a, "r").read(),
}
