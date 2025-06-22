package loader

import (
	"io/fs"
	"log"
	"path/filepath"

	"github.com/msp301/zb/util"
)

type Loader struct {
	Files []string
}

func New() *Loader {
	return &Loader{}
}

func (l *Loader) Load(paths ...string) error {
	for _, dir := range paths {
		err := filepath.WalkDir(dir, func(path string, d fs.DirEntry, err error) error {
			if err != nil {
				log.Fatalf(err.Error())
			}

			_, err = util.FileId(path)
			if err != nil {
				return nil
			}

			l.Files = append(l.Files, path)

			return nil
		})

		if err != nil {
			return err
		}
	}

	return nil
}
