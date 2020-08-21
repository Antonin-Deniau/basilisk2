#!/usr/bin/env python
import sys
import atexit
import os
import readline
import traceback

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
    sys.stdout.write(" ".join(map(display, e)))
    sys.stdout.write("\n")

def rep(e, env):
    b = read(e)
    c = evl(b, env)
    prnt(c)

repl_env = Env(None, [], [])
for k, v in ns.items():
    repl_env.set(k, v)

rep("(def! not (fn* (a) (if a false true)))", repl_env)
rep('(def! load-file (fn* (f) (eval (read-string (str "(do " (slurp f) "\\nnil)")))))', repl_env)

if len(sys.argv) >= 2:
    data = open(sys.argv[1], "r").readlines()
    for a in data:
        print(a.strip())

        try:
            rep(a, repl_env)
        except Exception as e:
            print(e)
            traceback.print_exc(file=sys.stdout)
    exit()

while True:
    try:
        rep(input("~>"), repl_env)
    except Exception as e:
        print(e)
        traceback.print_exc(file=sys.stdout)
