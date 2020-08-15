from parser import display, Name
from environment import Env

### SYMBOLS ###

def let_symbol(ast,env):
    if not isinstance(ast[1], tuple) and not isinstance(ast[1], list): raise Exception("Not a list or vector {}".format(ast[1]))
    if len(ast) != 3: raise Exception("Bad number or argument ({} for 3) for get* ({})".format(
        len(ast), display(ast)))

    new_env = Env(env)
    binding_list = ast[1]

    for i in zip(binding_list[::2], binding_list[1::2]):
        data = evl(i[1], new_env)
        new_env.set(i[0], data)

    return evl(ast[2], new_env)

def eval_function(ast, env):
    ev = eval_ast(ast, env)
    return ev[0](*ev[1::])

def def_symbol(ast, env):
    if not isinstance(ast[1], Name): raise Exception("Not a symbol {}".format(ast[1]))
    if len(ast) != 3: raise Exception("Bad number or argument ({} for 3) for def! ({})".format(
        len(ast), display(ast)))

    value = evl(ast[2], env)
    env.set(ast[1].name, value)
    return value

def core_switch(ast, env):
    if isinstance(ast[0], Name):
        if ast[0].name == "def!":
            return def_symbol(ast,env)
        if ast[0].name == "let*":
            return let_symbol(ast,env)

        return eval_function(ast, env)
    else:
        raise Exception("{} not callable".format(ast[0]))

### EVAL PART ###

def evl(ast, env):
    if isinstance(ast, tuple):
        if len(ast) == 0:
            return ast

        return core_switch(ast,env)

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

    return ast
