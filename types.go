package main


type BType interface { }

type BString struct {
	Meta *BType
	Value string
}

type BList struct {
	Meta *BType
	Value []*BType
}

type BBool struct {
	Meta *BType
	Value bool
}


type BNil struct {}
