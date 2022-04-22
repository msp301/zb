package main

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/msp301/zb/parser"
)

func main() {

	filename := os.Args[1]

	note := parser.Parse(filename)

	json, err := json.Marshal(note)
	if err != nil {
		panic(err)
	}

	fmt.Println(string(json))
}
