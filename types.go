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
