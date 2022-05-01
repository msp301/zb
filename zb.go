package main

import (
	"encoding/json"
	"fmt"
	"os"

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

	json, err := json.Marshal(book.Read())
	if err != nil {
		panic(err)
	}

	fmt.Println(string(json))
}

func isValidNote(note parser.Note, book notebook.Notebook) bool {
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
