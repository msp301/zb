package notebook

import (
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/msp301/zb/parser"
)

type Notebook struct {
	Path string
}

func New(path string) Notebook {
	return Notebook{
		Path: path,
	}
}

func (book *Notebook) Read() []parser.Note {
	var notes []parser.Note

	filepath.Walk(book.Path, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			log.Fatalf(err.Error())
		}

		if !strings.HasSuffix(info.Name(), ".md") {
			return nil
		}

		notes = append(notes, parser.Parse(path)...)

		return nil
	})

	return notes
}
