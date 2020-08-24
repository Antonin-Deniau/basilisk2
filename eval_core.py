from functools import reduce
import types

from parser import display
from basl_types import Fn, Name
from environment import Env

### SYMBOLS ###

def check_fn(ast):
    if len(ast) != 3:
        raise Exception("Bad number or argument ({} for 3) for fn* ({})".format(len(ast), display(ast)))

    if (Name("&") in ast[1]) and (ast[1].index(Name("&")) != len(ast[1]) - 2):
        raise Exception("Function should contain only one variadic argument")

def check_def(ast):
    if not isinstance(ast[1], Name):
        raise Exception("Not a symbol {}".format(ast[1]))
    if len(ast) != 3:
        raise Exception("Bad number or argument ({} for 3) for def! ({})".format(len(ast), display(ast)))

def check_let(ast):
    if not isinstance(ast[1], tuple) and not isinstance(ast[1], list):
        raise Exception("Not a list or vector {}".format(ast[1]))
    if len(ast) != 3:
        raise Exception("Bad number or argument ({} for 3) for get* ({})".format(len(ast), display(ast)))


def quasiquote_process_list(ast):
    res = tuple([])
    for elt in reversed(ast):
        if isinstance(elt, tuple) and len(elt) != 0 and isinstance(elt[0], Name) and elt[0].name == "splice-unquote":
            res = tuple([Name("concat"), elt[1], res])
        else:
            res = tuple([Name("cons"), quasiquote(elt), res])

    return res

def quasiquote(ast):
    if isinstance(ast, tuple):
        if len(ast) != 0 and isinstance(ast[0], Name):
            if isinstance(ast[0], Name) and ast[0].name == "unquote":
                return ast[1]
        return quasiquote_process_list(ast)

    if isinstance(ast, list):
        return tuple([Name("vec"), quasiquote_process_list(tuple(ast))])

    if isinstance(ast, dict) or isinstance(ast, Name):
        return tuple([Name("quote"), ast])

    return ast

### EVAL PART ###

def evl(ast, env):
    while True:
        if isinstance(ast, tuple):
            if len(ast) == 0: return ast

            if isinstance(ast[0], Name):
                if ast[0].name == "def!":
                    check_def(ast)
                    value = evl(ast[2], env)
                    return env.set(ast[1].name, value)

                if ast[0].name == "let*":
                    check_let(ast)
                    new_env = Env(env, [], [])
                    binding_list = ast[1]

                    for i in zip(binding_list[::2], binding_list[1::2]):
                        data = evl(i[1], new_env)
                        new_env.set(i[0], data)

                    ast, env = ast[2], new_env; continue

                if ast[0].name == "quote":
                    return ast[1]

                if ast[0].name == "quasiquoteexpand":
                    return quasiquote(ast[1])

                if ast[0].name == "quasiquote":
                    ast = quasiquote(ast[1]); continue

                if ast[0].name == "do":
                    res = None
                    for x in ast[1:-1]:
                        res = evl(x, env)
                    ast = ast[-1]; continue

                if ast[0].name == "if":
                    if len(ast) < 3: ast, env = None, env; continue
                    res_cond = evl(ast[1], env)

                    if type(res_cond) == bool and res_cond == True: ast = ast[2]; continue
                    if type(res_cond) == int: ast = ast[2]; continue
                    if type(res_cond) == float: ast = ast[2]; continue
                    if type(res_cond) == list: ast = ast[2]; continue
                    if type(res_cond) == tuple: ast = ast[2]; continue
                    if type(res_cond) == str: ast = ast[2]; continue

                    ast = ast[3] if len(ast) >= 4 else None; continue

                if ast[0].name == "fn*":
                    body = ast[2]
                    params = ast[1]

                    func = lambda *e: evl(body, Env(env, params, e))
                    return Fn(body, params, env, func)

            [f, *args] = eval_ast(ast, env)

            if isinstance(f, Fn):
                ast, env = f.ast, Env(f.env, f.params, args); continue

            if isinstance(f, types.LambdaType):
                return f(*args)

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

