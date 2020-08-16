class Env:
    def __init__(self, outer, binds, exprs):
        self.outer = outer
        self.vals = {}

        for i in zip(binds, exprs):
            self.set(i[0], i[1])

    def set(self, name, value):
        self.vals[name] = value

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
            raise Exception("Symbol not found {}".format(name))
