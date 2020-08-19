from parser import display

def prn(*a):
    print(" ".join([display(i) for i in a]))
    return None

def println(*a):
    print(" ".join([display(i, False) for i in a]))
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
    '=': lambda a, b: type(a) == type(b) and a == b,
    '<': lambda a, b: a < b,
    '<=': lambda a, b: a <= b,
    '>=': lambda a, b: a >= b,
    '>': lambda a, b: a > b,
    'pr-str': lambda *a: " ".join([display(i) for i in a]),
    'str': lambda *a: "".join([display(i, False) for i in a]),
    'prn': prn,
    'println': println,
}
