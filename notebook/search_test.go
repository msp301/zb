package notebook

import (
	"github.com/msp301/graph"
	"github.com/msp301/zb/parser"
	"reflect"
	"testing"
)

func TestSearch(t *testing.T) {
	g := graph.New()
	g.Add(1, "note", parser.Note{Content: "foo bar", Start: 1})
	g.Add(2, "note", parser.Note{Content: "# TICKET-123 title", Start: 1})
	book := &Notebook{Notes: g}

	got := book.Search("ticket-123")
	want := []Result{{Context: "# TICKET-123 title", Line: 1, Value: graph.Vertex{Id: 2, Label: "note", Value: parser.Note{Content: "# TICKET-123 title", Start: 1}}}}

	if !reflect.DeepEqual(got, want) {
		t.Fatalf("Expected: %+v\nGot: %+v", want, got)
	}
}
