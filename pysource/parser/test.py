

macros = {}


def set_macro_character(char, func):
    macros[char] = func


def read(stream, eof_error, eof_value, recursive):
    res = eof_value
    no_value = True

    stream_pointer = 0
    last_sym = 0
    while x <= len(stream):
        stream_pointer += 1

        if stream[last_sym:stream_pointer] in macros:
            ast = macro(stream, stream[last_sym:stream_pointer])
            res = eval(ast)
            no_value = False

    if eof_error != None and no_value:
        return BaslException()

    return res


