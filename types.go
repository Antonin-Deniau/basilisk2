package main

import (
	"github.com/alecthomas/participle/lexer"
)

// PARSABLE

type BType struct {
	Pos lexer.Position
	
	BNumber        *BNumber        `@@`
	BKeyword       *BKeyword       `| @@`
	BString        *BString        `| @@`
	BList          *BList          `| @@`
	BMetadata      *BMetadata      `| @@`
	BDeref         *BDeref         `| @@`
	BHashMap       *BHashMap       `| @@`
	BVector        *BVector        `| @@`
	BQuote         *BQuote         `| @@`
	BQuasiquote    *BQuasiquote    `| @@`
	BSpliceUnquote *BSpliceUnquote `| @@`
	BUnquote       *BUnquote       `| @@`
	BName          *BName          `| @@`
	BBoolean       *BBoolean       `| @@`
	BVariadic      *BVariadic      `| "&"`
	BNil           *BNil           `| @Nil`
}

type BNil struct {
	Pos lexer.Position
}

type BVariadic struct {
	Pos lexer.Position
}

type BBoolean struct {
	Pos lexer.Position

	Value bool `@Bool`
}

type BList struct {
	Pos lexer.Position

	Values []*BType `"(" @@* ")"`
}

type BMetadata struct {
	Metadata *BType `"^" @@`
	Value    *BType ` @@`
}

type BDeref struct {
	Pos lexer.Position

	Value *BType `"@" @@`
}

type BHashMap struct {
	Pos lexer.Position

	Map []*BHashMapEntry `"{" @@* "}"`
}

type BHashMapEntry struct {
	Pos lexer.Position

	Key   *BType `@@`
	Value *BType ` @@`
}

type BVector struct {
	Pos lexer.Position

	Values []*BType `"[" @@* "]"`
}

type BQuote struct {
	Pos lexer.Position

	Value *BType `"'" @@`
}

type BQuasiquote struct {
	Pos lexer.Position

	Value *BType "\"`\" @@"
}

type BSpliceUnquote struct {
	Pos lexer.Position

	Value *BType `"~@" @@`
}

type BUnquote struct {
	Pos lexer.Position

	Value *BType `"~" @@`
}

type BString struct {
	Pos lexer.Position

	Value string `@String`
}

func (t *BString) Capture(s []string) error {
	t.Value = Unescape(s[0])
	return nil
}

type BNumber struct {
	Pos lexer.Position

	Value float64 `@Number`
}

type BName struct {
	Pos lexer.Position

	Value string `@Ident`
}

type BKeyword struct {
	Pos lexer.Position

	Value string `":" @Ident`
}

// NON PARSABLE

/*
type BFn struct {
	Ast     *BType
	Params  []*BType
	Env     *Env
	Fn      *BFn
	IsMacro bool
	Meta    *BType
}

func NewBFn() {
    def __init__(self, ast, params, env, fn, is_macro=False, meta=None):
        self.ast = ast
        self.params = params
        self.env = env
        self.fn = fn
        self.is_macro = is_macro
        self.meta = None
}
*/

type BAtom struct {
	Value *BType
}

/*
func NewBAtom() {
    def __init__(self, data):
        self.data = data
}

func (BAtom) Reset() {
    def reset(self, a):
        self.data = a
        return a
}

type BException {
	IsRaw bool
	Message *BType
}

func NewBException() {
    def __init__(self, parent):
        self.is_raw = isinstance(parent, Exception)
        self.message = str(parent) if self.is_raw else parent
}
*/
