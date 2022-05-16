package main

import (
	"encoding/json"
	"fmt"
	"os"
	"strconv"
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
	case "note":
		id, _ := strconv.ParseUint(os.Args[3], 0, 64)
		render(book.SearchRelated(id))
	case "tag":
		render(book.SearchByTag(os.Args[3]))
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

func render(vertices []graph.Vertex) {
	for _, vertex := range vertices {
		switch val := vertex.Properties["Value"].(type) {
		case parser.Note:
			fmt.Printf("%s - %s\n", val.File, val.Title)
		case string:
			fmt.Printf("%s\n", val)
		}
		//jsonStr, err := json.Marshal(vertex)
		//if err != nil {
		//	panic("Failed to encode")
		//}
		//fmt.Println(string(jsonStr))
	}
}
