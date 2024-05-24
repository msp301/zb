package cmd

import (
	"fmt"
	"log"

	"github.com/msp301/zb/bookshelf"
	"github.com/msp301/zb/graph"
	"github.com/msp301/zb/notebook"
	"github.com/msp301/zb/pager"
	"github.com/msp301/zb/parser"
	"github.com/spf13/viper"
)

func bookDirs() []string {
	return viper.GetStringSlice("directory")
}

func book() *notebook.Notebook {
	err, book := bookshelf.Read(bookDirs())
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
	pager, err := pager.Open()
	if err != nil {
		log.Fatal(err)
	}

	go func() {
		defer pager.Close()

		for _, result := range results {
			switch val := result.Value.(type) {
			case graph.Vertex:
				switch vertex := val.Properties["Value"].(type) {
				case parser.Note:
					if len(result.Context) > 0 {
						pager.Writef("%s:%d: - %s\n", vertex.File, result.Line, result.Context)
					} else {
						pager.Writef("%s:%d: - %s\n", vertex.File, vertex.Start, vertex.Title)
					}
				case string:
					pager.Writef("%s\n", vertex)
				}
			}
		}
	}()

	if err := pager.Wait(); err != nil {
		log.Fatal(err)
	}
}
