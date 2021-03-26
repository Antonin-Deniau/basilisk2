class Name:
    def __init__(self, name, file=None, line=None, col=None):
        self.name = name
        self.file = file
        self.line = line
        self.col = col

    def __hash__(self):
        return hash(self.name)

    def __repr__(self):
        return self.name

    def __eq__(self, a):
        return self.name == a

class Keyword:
    def __init__(self, name):
        self.name = name

    def __hash__(self):
        return hash(self.name)

    def __repr__(self):
        return self.name

    def __eq__(self, a):
        return self.name == a

class Fn:
    def __init__(self, ast, params, env, fn, is_macro=False, meta=None):
        self.ast = ast
        self.params = params
        self.env = env
        self.fn = fn
        self.is_macro = is_macro
        self.meta = meta

class Atom:
    def __init__(self, data):
        self.data = data

    def reset(self, a):
        self.data = a
        return a

class BaslException(Exception):
    def __init__(self, parent, stack=None):
        self.is_raw = isinstance(parent, Exception)
        m = str(parent) if self.is_raw else parent

        if stack != None:
            self.message = "StackTrace: \n{}\n\nException: {}".format("\n".join(stack), m)
        else:
            self.message = m

    def __str__(self):
        return self.message
