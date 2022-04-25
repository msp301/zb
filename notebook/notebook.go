package notebook

import (
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/msp301/zb/parser"
)

type Notebook struct {
	Filters []FilterFunc
	lookup  map[uint64]parser.Note
	Path    string
	Notes   []parser.Note
}

type FilterFunc func(note parser.Note) bool

func New(path string) Notebook {
	return Notebook{
		Path:   path,
		lookup: map[uint64]parser.Note{},
	}
}

func (book *Notebook) AddFilter(filter FilterFunc) {
	book.Filters = append(book.Filters, filter)
}

func (book *Notebook) Read() []parser.Note {
	var filteredNotes []parser.Note

	filepath.Walk(book.Path, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			log.Fatalf(err.Error())
		}

		if !strings.HasSuffix(info.Name(), ".md") {
			return nil
		}

		fileNotes := parser.Parse(path)
		book.Notes = append(book.Notes, fileNotes...)

	NOTE:
		for _, note := range fileNotes {
			book.lookup[note.Id] = note

			for _, filter := range book.Filters {
				if filter(note) {
					continue NOTE
				}
			}
			filteredNotes = append(filteredNotes, note)
		}

		return nil
	})

	return filteredNotes
}

func (book *Notebook) GetNote(noteId uint64) (parser.Note, bool) {
	note, ok := book.lookup[noteId]
	return note, ok
}
