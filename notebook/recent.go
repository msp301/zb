package notebook

import (
	"sort"

	"github.com/msp301/graph"
	"github.com/msp301/zb"
)

func (book *Notebook) Recent(limit int) []Result {
	var notes []graph.Vertex
	traversal := graph.Traversal(book.Notes)
	for vertex := range traversal.V().HasLabel("note").Iterate() {
		if _, ok := vertex.Value.(zb.Note); ok {
			notes = append(notes, vertex)
		}
	}

	sort.Slice(notes, func(i, j int) bool {
		return notes[i].Id > notes[j].Id
	})

	if len(notes) > limit {
		notes = notes[:limit]
	}

	results := make([]Result, 0, len(notes))
	for _, vertex := range notes {
		results = append(results, Result{
			Line:  vertex.Value.(zb.Note).Start,
			Value: vertex,
		})
	}

	return results
}
