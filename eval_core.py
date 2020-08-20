import types

from parser import display, Name
from environment import Env

### SYMBOLS ###

def fn_symbol(ast, env):
    if len(ast) != 3: raise Exception("Bad number or argument ({} for 3) for fn* ({})".format(
        len(ast), display(ast)))

    if (Name("&") in ast[1]) and (ast[1].index(Name("&")) != len(ast[1]) - 2):
        raise Exception("Function should contain only one variadic argument")

    return lambda *e: evl(ast[2], Env(env, ast[1], e)), env

def if_symbol(ast, env):
    if len(ast) < 3: return None
    res_cond = evl(ast[1], env)

    if type(res_cond) == bool and res_cond == True: return ast[2], env
    if type(res_cond) == int: return ast[2], env
    if type(res_cond) == float: return ast[2], env
    if type(res_cond) == list: return ast[2], env
    if type(res_cond) == tuple: return ast[2], env
    if type(res_cond) == str: return ast[2], env

    return ast[3], env if len(ast) >= 4 else None, env

def do_symbol(ast, env):
    res = None
    for x in ast[1:-1]:
        res = evl(x, env)
    return ast[-1], env

def let_symbol(ast,env):
    if not isinstance(ast[1], tuple) and not isinstance(ast[1], list): raise Exception("Not a list or vector {}".format(ast[1]))
    if len(ast) != 3: raise Exception("Bad number or argument ({} for 3) for get* ({})".format(
        len(ast), display(ast)))

    new_env = Env(env, [], [])
    binding_list = ast[1]

    for i in zip(binding_list[::2], binding_list[1::2]):
        data = evl(i[1], new_env)
        new_env.set(i[0], data)

    return ast[2], new_env

def eval_function(ast, env):
    ev = eval_ast(ast, env)
    print(ev)
    return ev[0](*ev[1::]), env

def def_symbol(ast, env):
    if not isinstance(ast[1], Name): raise Exception("Not a symbol {}".format(ast[1]))
    if len(ast) != 3: raise Exception("Bad number or argument ({} for 3) for def! ({})".format(
        len(ast), display(ast)))

    value = evl(ast[2], env)
    env.set(ast[1].name, value)
    return value, env

### EVAL PART ###

def evl(ast, env):
    while True:
        if isinstance(ast, tuple):
            if len(ast) == 0: return ast

            if isinstance(ast[0], Name):
                if ast[0].name == "def!": ast, env = def_symbol(ast,env); continue
                if ast[0].name == "let*": ast, env = let_symbol(ast,env); continue
                if ast[0].name == "do":   ast, env = do_symbol(ast,env); continue
                if ast[0].name == "if":   ast, env = if_symbol(ast,env); continue
                if ast[0].name == "fn*":  ast, env = fn_symbol(ast,env); continue

            ast, env = eval_function(ast, env); continue

        return eval_ast(ast, env)

def eval_ast(ast, env):
    if isinstance(ast, dict):
        return { k: evl(v, env) for k,v in ast.items() }

    if isinstance(ast, Name):
        return env.get(ast.name)

    if isinstance(ast, list):
        return list([evl(a, env) for a in ast])

    if isinstance(ast, tuple):
        return tuple([evl(x, env) for x in ast])

    return ast, env
