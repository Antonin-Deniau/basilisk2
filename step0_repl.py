import sys, parser
import atexit
import os
import readline
histfile = os.path.join(os.path.expanduser("~"), ".python_history")

try:
    readline.read_history_file(histfile)
    h_len = readline.get_current_history_length()
except FileNotFoundError:
    open(histfile, 'wb').close()
    h_len = 0

def save(prev_h_len, histfile):
    new_h_len = readline.get_current_history_length()
    readline.set_history_length(1000)
    readline.append_history_file(new_h_len - prev_h_len, histfile)
atexit.register(save, h_len, histfile)

def read():
    e = input("user>")
    return parser.parse(e)

def format_type(data):
    if isinstance(data, list):
        return [format_type(d) for d in data]
    if isinstance(data, float):
        return { "type": "number", "value": data }
    if isinstance(data, int):
        return { "type": "number", "value": data }
    if isinstance(data, dict):
        return { "type": "hashmap", "value": { format_type(k): format_type(v) for k,v in data.items()} }
    if isinstance(data, str):
        return { "type": "string", "value": data }

def evl(ast, env):
    if isinstance(ast, list):
        if len(ast) == 0:
            return ast
        else:
            ev = eval_ast(ast, env)
            return ev[0](*ev[1::])
    else:
        return eval_ast(ast, env)

def prnt(e):
    sys.stdout.write(parser.display(e))
    sys.stdout.write("\n")

def rep(env):
    b = read()
    c = [evl(d, env) for d in b]
    prnt(c)

def eval_ast(ast, env):
    if isinstance(ast, dict):
        if ast["type"] == "number":
            return ast["value"]
        if ast["type"] == "string":
            return ast["value"]
        if ast["type"] == "vector":
            return [evl(a, env) for a in ast["value"]]
        if ast["type"] == "hashmap":
            return { k[0]["value"]: evl(k[1], env) for k in ast["value"]}
        if ast["type"] == "name":
            if ast["value"] in env:
                return env[ast["value"]]
            else:
                raise Exception("Symbol not found: {}".format(ast["value"]))

    if isinstance(ast, list):
        return [evl(x, env) for x in ast]

    return ast

repl_env = {
        '+': lambda a,b: a+b,
        '-': lambda a,b: a-b,
        '*': lambda a,b: a*b,
        '/': lambda a,b: int(a/b),
}

while True:
    rep(repl_env)
