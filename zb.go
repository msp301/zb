package main

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/msp301/zb/graph"
	"github.com/msp301/zb/notebook"
	"github.com/msp301/zb/parser"
)

func main() {

	var action, dirname string

	if len(os.Args) == 2 {
		action = "parse"
		dirname = os.Args[1]
	} else {
		action = os.Args[1]
		dirname = os.Args[2]
	}

	book := notebook.New(dirname)

	switch action {
	case "check":
		book.AddFilter(func(note parser.Note) bool {
			return isValidNote(note, book)
		})
	}

	notes := book.Read()

	switch action {
	case "outline":
		book.Notes.Walk(func(vertex graph.Vertex, depth int) bool {
			indent := strings.Repeat("\t", depth)
			switch val := vertex.Properties["Value"].(type) {
			case parser.Note:
				fmt.Printf("%s%s - %s\n", indent, val.File, val.Title)
			}
			return true
		})
	case "tags":
		for _, tag := range book.Tags() {
			fmt.Println(tag)
		}
	default:
		jsonStr, err := json.Marshal(notes)
		if err != nil {
			panic(err)
		}
		fmt.Println(string(jsonStr))
	}
}

func isValidNote(note parser.Note, book *notebook.Notebook) bool {
	if note.Id == 0 {
		return false
	}

	for _, link := range note.Links {
		if !book.IsNote(link) {
			return false
		}
	}

	return true
}
