package cmd

import (
	"fmt"
	"github.com/msp301/zb/bookshelf"
	"github.com/msp301/zb/graph"
	"github.com/msp301/zb/notebook"
	"github.com/msp301/zb/parser"
	"github.com/spf13/viper"
	"log"
)

func book() *notebook.Notebook {
	dirs := viper.GetStringSlice("directory")
	err, book := bookshelf.Read(dirs)
	if err != nil {
		log.Fatalf("Failed to read notebook: %s", err)
	}
	return book
}

func render(vertices []graph.Vertex) {
	for _, vertex := range vertices {
		switch val := vertex.Properties["Value"].(type) {
		case parser.Note:
			fmt.Printf("%s - %s\n", val.File, val.Title)
		case string:
			fmt.Printf("%s\n", val)
		}
	}
}

func renderResults(results []notebook.Result) {
	for _, result := range results {
		switch val := result.Value.(type) {
		case graph.Vertex:
			switch vertex := val.Properties["Value"].(type) {
			case parser.Note:
				if len(result.Context) > 0 {
					fmt.Printf("%s - %s\n", vertex.File, result.Context)
				} else {
					fmt.Printf("%s - %s\n", vertex.File, vertex.Title)
				}
			case string:
				fmt.Printf("%s\n", vertex)
			}
		}
	}
}
