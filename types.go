package main


type BType interface { }

type BString struct {
	Meta *BType
	Value string
}

type BInt struct {
	Meta *BType
	Value int64
}

type BVariadic struct {}

type BHashmap struct {
	Meta *BType
	Value map[*BType]*BType
}

type BList struct {
	Meta *BType
	Value []*BType
}

type BBool struct {
	Meta *BType
	Value bool
}

type BName struct {
	Meta *BType
	Value string
}

type BKeyword struct {
	Meta *BType
	Value string
}


type BNil struct {}
