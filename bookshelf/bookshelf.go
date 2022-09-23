package bookshelf

import "github.com/msp301/zb/notebook"

func Read(paths []string) (error, *notebook.Notebook) {
	book := notebook.New()
	for _, path := range paths {
		err := book.Load(path)
		if err != nil {
			return err, nil
		}
	}
	book.Read()
	return nil, book
}
