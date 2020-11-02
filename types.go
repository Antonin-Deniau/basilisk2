package main

import (
	"github.com/alecthomas/participle/lexer"
)

type BType struct {
	Pos lexer.Position
	
	Index int

	Line int `@Number`

	Number        *float64    `  @Number`
	Variable      *string     `| @Ident`
	String        *string     `| @String`
	Call          *Call       `| @@`
	Subexpression *Expression `| "(" @@ ")"`
}

type BList struct {
}

type BMetadata struct {
}

type BDeref struct {
}

type BHashMap struct {
}

type BVector struct {
}

type BQuote struct {
}

type BQuasiquote struct {
}

type BSpliceUnquote struct {
}

type BUnquote struct {
}

type BString struct {
}

type BNumber struct {
}

type BName struct {
	Pos lexer.Position

	name string `@Ident`
}

type BKeyword struct {
    def __init__(self, name):
        self.name = name

    def __hash__(self):
        return hash(self.name)

    def __repr__(self):
        return self.name

    def __eq__(self, a):
        return self.name == a
}

type BFn struct {
    def __init__(self, ast, params, env, fn, is_macro=False, meta=None):
        self.ast = ast
        self.params = params
        self.env = env
        self.fn = fn
        self.is_macro = is_macro
        self.meta = None
}

type BAtom struct {
    def __init__(self, data):
        self.data = data

    def reset(self, a):
        self.data = a
        return a
}

type BException {
    def __init__(self, parent):
        self.is_raw = isinstance(parent, Exception)
        self.message = str(parent) if self.is_raw else parent
}
