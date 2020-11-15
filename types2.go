package main

// PARSABLE

type BType {
	Pos Position
}

type Position struct {
	Line int
	Col int
}

type BNil struct {
	BType

	Value bool
}

type BVariadic struct {
	BType
}

type BBoolean struct {
	BType

	Value bool
}

type BList struct {
	BType

	Values []*BType
}

type BMetadata struct {
	BType

	Metadata *BType
	Value    *BType
}

type BDeref struct {
	BType

	Value *BType
}

type BHashMap struct {
	BType

	Map []*BHashMapEntry
}

type BHashMapEntry struct {
	Key   *BType
	Value *BType
}

type BVector struct {
	BType

	Values []*BType
}

type BQuote struct {
	BType

	Value *BType
}

type BQuasiquote struct {
	Pos Position

	Value *BType
}

type BSpliceUnquote struct {
	Pos Position

	Value *BType
}

type BUnquote struct {
	Pos Position

	Value *BType
}

type BString struct {
	Pos Position

	Value string
}

type BNumber struct {
	Pos Position

	Value float64
}

type BName struct {
	Pos Position

	Value string
}

type BKeyword struct {
	Pos Position

	Value string
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
