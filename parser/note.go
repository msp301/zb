package parser

type Note struct {
	File  string
	Id    uint64
	Links []uint64
	Tags  []string
	Title string
}
