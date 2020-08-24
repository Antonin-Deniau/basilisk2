#!/usr/bin/env python
import sys
import atexit
import json
import os
import readline
import traceback
from lark import UnexpectedInput

from core import ns
from parser import parse, display
from basl_types import Name
from eval_core import evl
from environment import Env

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




def read(e):
    return parse(e)

def prnt(e):
    if e != None: sys.stdout.write(display(e))
    sys.stdout.write("\n")

def rep(e, env):
    b = read(e)
    c = evl(b, env)
    prnt(c)

def load_str(e, env):
    b = read(e)
    evl(b, env)



repl_env = Env(None, [], [])
for k, v in ns.items():
    repl_env.set(k, v)

repl_env.set("eval", lambda e: evl(e, repl_env))
repl_env.set("*ARGV*", tuple(sys.argv))

load_str("(def! not (fn* (a) (if a false true)))", repl_env)
load_str('(def! load-file (fn* (f) (eval (read-string (str "(do " (slurp f) "\\nnil)")))))', repl_env)

if len(sys.argv) >= 2:
    if sys.argv[1] == "test":
        data = open(sys.argv[2], "r").readlines()
        for a in data:
            if a.strip == "": continue
            print(a.strip())

            try:
                rep(a, repl_env)
            except UnexpectedInput as e:
                print(e)
                traceback.print_exc(file=sys.stdout)
    else:
        load_str("(load-file " + json.dumps(sys.argv[1]) + ")", repl_env)
else:
    while True:
        try:
            rep(input("~>"), repl_env)
        except Exception as e:
            print(e)
            traceback.print_exc(file=sys.stdout)
