#!/usr/bin/env python
import sys, parser
import atexit
import os
import readline
import traceback

histfile = os.path.join(os.path.expanduser("~"), ".basilisk_history")

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
    e = input("~>")
    return parser.parse(e)

def evl(ast, env):
    if isinstance(ast, tuple):
        if len(ast) == 0:
            return ast
        else:
            ev = eval_ast(ast, env)
            return ev[0](*ev[1::])

    return eval_ast(ast, env)

def prnt(e):
    sys.stdout.write(" ".join([parser.display(s) for s in e]))
    sys.stdout.write("\n")

def rep(env):
    b = read()
    c = [evl(d, env) for d in b]
    prnt(c)

def eval_ast(ast, env):
    if isinstance(ast, dict):
        return { k: evl(v, env) for k,v in ast.items() }

    if isinstance(ast, parser.Name):
        if ast.name in env:
            return env[ast.name]
        else:
            raise Exception("Symbol not found: {}".format(ast.name))

    if isinstance(ast, list):
        return [evl(a, env) for a in ast]

    if isinstance(ast, tuple):
        return tuple(evl(x, env) for x in ast)

    return ast

repl_env = {
        '+': lambda a,b: a+b,
        '-': lambda a,b: a-b,
        '*': lambda a,b: a*b,
        '/': lambda a,b: int(a/b),
}

while True:
    try:
        rep(repl_env)
    except Exception as e:
        print(e)
        traceback.print_exc(file=sys.stdout)
