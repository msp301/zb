package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/msp301/zb/parser"
)

func main() {

	dirname := os.Args[1]

	filepath.Walk(dirname, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			log.Fatalf(err.Error())
		}

		if !strings.HasSuffix(info.Name(), ".md") {
			return nil
		}

		note := parser.Parse(path)

		json, err := json.Marshal(note)
		if err != nil {
			panic(err)
		}

		fmt.Println(string(json))
		return nil
	})
}
