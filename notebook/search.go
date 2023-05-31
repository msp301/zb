package notebook

import (
	"github.com/msp301/zb/graph"
	"github.com/msp301/zb/parser"
	"github.com/msp301/zb/util"
)

func (book *Notebook) Search(query string) []Result {
	results := []Result{}
	traversal := graph.Traversal(book.Notes)
	for vertex := range traversal.V().HasLabel("note").Iterate() {
		context := []string{""}
		switch val := vertex.Properties["Value"].(type) {
		case parser.Note:
			matched, ok := util.Context(val.Content, query)
			if ok {
				context = matched
			} else {
				continue
			}
		}

		for _, context := range context {
			result := Result{
				Context: context,
				Value:   vertex,
			}
			results = append(results, result)
		}
	}
	return results
}
