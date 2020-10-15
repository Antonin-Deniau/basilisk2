def call(name, *args):
    attrs = name.split(".")
    mod = attrs[-1]
    path = attrs[1:-1]

    return reduce(lambda acc, arr: getattr(acc, arr), attrs, __import__(mod))(*args)

