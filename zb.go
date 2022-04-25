package main

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/msp301/zb/notebook"
)

func main() {

	dirname := os.Args[1]

	book := notebook.New(dirname)

	json, err := json.Marshal(book.Read())
	if err != nil {
		panic(err)
	}

	fmt.Println(string(json))
}
