from basl_types import Name, BaslException

class Env:
    def __init__(self, outer, binds, exprs, new_trace=None):
        self.outer = outer
        self.vals = {}
        if outer != None:
          self.stack = [*outer.stack, new_trace] if new_trace != None else outer.stack
        else:
          self.stack = []

        if len(binds) != len(exprs):
            if Name("&") not in binds:
                raise BaslException("Function should contain {} parametter".format(len(binds)), self.stack)

            if len(exprs) < len(binds) - 2:
                raise BaslException("Function should contain at least {} parametter".format(len(binds) - 2), self.stack)

        for i in zip(binds, exprs):
            if Name("&") == i[0]: break
            self.set(i[0], i[1])

        if Name("&") in binds:
            self.set(binds[-1], tuple(exprs[len(binds) - 2::]))

    def set(self, name, value):
        self.vals[name] = value
        return value

    def find(self, name):
        if name in self.vals:
            return self
        else:
            if self.outer is not None:
                return self.outer.find(name)
            else:
                return None

    def get(self, name):
        env = self.find(name)

        if env is not None:
            return env.vals[name]
        else:
            raise BaslException("{} not found".format("'{}'".format(name)), self.stack)
