class Name:
    def __init__(self, name):
        self.name = name

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
    def __init__(self, ast, params, env, fn):
        self.ast = ast
        self.params = params
        self.env = env
        self.fn = fn
