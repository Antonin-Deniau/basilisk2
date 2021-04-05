from functools import reduce
import base64, time, types, re, sys, io

from parser_t import display, parse
from lark import UnexpectedInput, UnexpectedToken
from basl_types import Name, Atom, Fn, BaslException, Keyword

def with_meta(e, meta):
    return e

def readline(e):
    try:
        return input(e)
    except BaseException:
        return None

def read_string(a):
    try: 
        return parse(a)
    except UnexpectedToken as e:
        raise BaslException("EOF/Parsing error around: \n" + e.get_context(a, 200))
    except UnexpectedInput as e:
        raise BaslException("EOF/Parsing error around: \n" + e.get_context(a, 200))
    except Exception as e:
        raise BaslException(e)

def regex_match(reg, text):
    e = re.match(reg, text)
    if e != None:
        return [e.group(0), list(e.groups())]
    else:
        return None

def prn(*a):
    print(" ".join([display(i, True) for i in a]))
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

    if type(a) == dict and type(b) == dict:
        for k in a.keys():
            if k not in b: return False

            bv = b[k]
            av = a[k]
            if not equality(bv, av): return False
        return True

    return type(a) == type(b) and a == b

def swap(a, b, *c):
    if isinstance(b, Fn):
        return a.reset(b.fn(a.data, *c))
    else:
        return a.reset(b(a.data, *c))

def ret_func(a):
    if isinstance(a, Fn):
        return a.fn
    else:
        return a

def appl(a, *b):
    is_arr = lambda x: isinstance(x, tuple) or isinstance(x, list)
    args = reduce(lambda acc, arr: [*acc, *arr] if is_arr(arr) else [*acc, arr], b, [])
    return ret_func(a)(*args)

def basl_map(a, b):
    s = map(ret_func(a), b)
    return tuple(s)

def basl_reduce(a, b):
    return reduce(ret_func(a), b)

def pr_str(*a):
    return " ".join([display(i, True) for i in a])

def peek(f):
    pos = f.tell()
    line = f.read(1)
    f.seek(pos)
    return line

ns = {
    # STREAMS FUNCS
    'output-stream': lambda e: open(e, 'r'),
    'input-stream': lambda e: open(e, 'w'),
    'string-stream': lambda e: io.StringIO(e),
    '*stdin*': sys.stdin,
    '*stdout*': sys.stdout,
    'read-all': lambda e: e.read(),
    'read-byte': lambda e: e.read(1),
    'peek-byte': peek,
    'write-byte': lambda e, b: e.write(b),
    # https://dept-info.labri.fr/~strandh/Teaching/MTP/Common/David-Lamkins/chapter19.html
    
    '&&': lambda a,b: a and b,
    '||': lambda a,b: a or b,
    'ord': lambda a: ord(a),
    '%': lambda a,b: a % b,
    'chr': lambda a: chr(a),
    'subs': lambda a, b, *c: type(a)(a[b:c[0] if len(c) != 0 else None]),
    '+': lambda a,b: a+b,
    '-': lambda a,b: a-b,
    '*': lambda a,b: a*b,
    '**': lambda a,b: a**b,
    '/': lambda a,b: a/b,
    'list': lambda *a: tuple(a),
    'list?': lambda a: isinstance(a, tuple),
    'empty?': lambda a: len(a) == 0,
    'count': lambda a: 0 if a == None else len(a),
    '=': equality,
    '<': lambda a, b: a < b,
    '<=': lambda a, b: a <= b,
    '>=': lambda a, b: a >= b,
    '>': lambda a, b: a > b,
    'pr-str': pr_str,
    'str': lambda *a: "".join([display(i, False) for i in a]),
    'prn': prn,
    'println': println,
    'read-string': read_string,
    #'slurp': lambda a: open(a, "r").read(),
    'atom': lambda a: Atom(a),
    'atom?': lambda a: isinstance(a, Atom),
    'deref': lambda a: a.data if isinstance(a, Atom) else nil,
    'reset!': lambda a, b: a.reset(b) if isinstance(a, Atom) else nil,
    'swap!': swap,
    'cons': lambda a, b: tuple([a,*b]),
    'conj': lambda a, b: [*a,b],
    'concat': lambda *a: tuple(reduce(lambda acc, arr: [*acc, *arr], a, ())),
    'vec': lambda a: list(a) if isinstance(a, list) or isinstance(a, tuple) else [a],
    'nth': lambda a, i: a[i],
    'first': lambda a: a[0] if a != None and len(a) != 0 else None,
    'rest': lambda a: tuple(a[1:]) if a != None and len(a) != 0 else tuple(),
    'apply': appl,
    'map': basl_map,
    'reduce': basl_reduce,
    'nil?': lambda e: isinstance(e, type(None)),
    'true?': lambda e: e == True,
    'false?': lambda e: e == False,
    'symbol?': lambda e: isinstance(e, Name),
    'symbol': lambda e: Name(e),
    'keyword': lambda e: Keyword(e.name if isinstance(e, Keyword) else e),
    'keyword?': lambda e: isinstance(e, Keyword),
    'vector': lambda *e: list([*e]),
    'vector?': lambda e: isinstance(e, list),
    'sequential?': lambda e: isinstance(e, list) or isinstance(e, tuple),
    'hash-map': lambda *e: { k[0]: k[1] for k in zip(e[::2], e[1::2]) },
    'map?': lambda e: isinstance(e, dict),
    'assoc': lambda a, *e: {**a, **{ k[0]: k[1] for k in zip(e[::2], e[1::2]) } },
    'dissoc': lambda e, *di: {key: val for key, val in e.items() if key not in di },
    'get': lambda e, k: e.get(k, None) if e is not None else None,
    'contains?': lambda e, k: k in e,
    'keys': lambda e: tuple(e.keys()),
    'vals': lambda e: tuple(e.values()),
    'time-ms': lambda: int(round(time.time() * 1000)),
    'number?': lambda e: (isinstance(e, float) or isinstance(e, int)) and not isinstance(e, bool),
    'string?': lambda e: isinstance(e, str),
    'fn?': lambda e: (isinstance(e, Fn) and e.is_macro == False) or isinstance(e, types.LambdaType),
    'macro?': lambda e: isinstance(e, Fn) and e.is_macro,
    'conj': lambda e, *rest: tuple([*reversed(rest), *e]) if isinstance(e, tuple) else list([*e, *rest]),
    'quux-eval': lambda e: eval(e),
    'readline': readline,
    'meta': lambda e: e.meta if hasattr(e, "meta") else None,
    'with-meta': with_meta,
    'seq': lambda e: tuple(e) if e != None and len(e) != 0 else None,
    'regex-match': regex_match,
}
