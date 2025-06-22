package bookshelf

import (
	"github.com/msp301/zb/loader"
	"github.com/msp301/zb/notebook"
)

func Read(paths []string) (error, *notebook.Notebook) {
	loader := loader.New()
	err := loader.Load(paths...)
	if err != nil {
		return err, nil
	}

	book := notebook.New(loader.Files)

	return nil, book
}
