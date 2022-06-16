package cmd

import (
	"fmt"
	"github.com/msp301/zb/graph"
	"github.com/msp301/zb/notebook"
	"github.com/msp301/zb/parser"
	"github.com/spf13/viper"
)

func book() *notebook.Notebook {
	// TODO: Add support for reading multiple notebook directories
	dirs := viper.GetStringSlice("directory")
	book := notebook.New(dirs[0])
	book.Read()
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
			}
		case string:
			fmt.Printf("%s\n", val)
		}
	}
}
