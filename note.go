package zb

type Note struct {
	Content string
	Start   int
	File    string
	Id      uint64
	Links   []uint64
	Tags    []string
	Title   string
}
