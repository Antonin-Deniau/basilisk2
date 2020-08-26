from functools import reduce

from lark import UnexpectedInput
from parser import display, parse
from basl_types import Name, Atom, Fn, BaslException

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
    except UnexpectedInput as e:
        raise Exception("Erreur dans la chaine, par la: \n" + e.get_context(a, 200))
    except IndexError:
        return None

def swap(a, b, *c):
    if isinstance(b, Fn):
        return a.reset(b.fn(a.data, *c))
    else:
        return a.reset(b(a.data, *c))

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
    'atom': lambda a: Atom(a),
    'atom?': lambda a: isinstance(a, Atom),
    'deref': lambda a: a.data if isinstance(a, Atom) else nil,
    'reset!': lambda a, b: a.reset(b) if isinstance(a, Atom) else nil,
    'swap!': swap,
    'cons': lambda a, b: tuple([a,*b]),
    'concat': lambda *a: tuple(reduce(lambda acc, arr: [*acc, *arr], a, ())),
    'vec': lambda a: list(a) if isinstance(a, list) or isinstance(a, tuple) else [a],
    'nth': lambda a, i: a[i],
    'first': lambda a: a[0] if a != None and len(a) != 0 else None,
    'rest': lambda a: tuple(a[1:]) if a != None and len(a) != 0 else tuple(),
    'throw': lambda a: raise BaslException(a),
}
