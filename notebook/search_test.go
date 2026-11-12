package notebook

import (
	"github.com/msp301/graph"
	"github.com/msp301/zb"
	"reflect"
	"testing"
)

func TestSearch(t *testing.T) {
	g := graph.New()
	g.Add(1, "note", zb.Note{Content: "foo bar", Start: 1})
	g.Add(2, "note", zb.Note{Content: "# TICKET-123 title", Start: 1})
	book := &Notebook{Notes: g}

	got := book.Search("ticket-123")
	want := []Result{{Context: "# TICKET-123 title", Line: 1, Value: graph.Vertex{Id: 2, Label: "note", Value: zb.Note{Content: "# TICKET-123 title", Start: 1}}}}

	if !reflect.DeepEqual(got, want) {
		t.Fatalf("Expected: %+v\nGot: %+v", want, got)
	}
}

func TestSearchReportsLineOfMatchInLaterParagraph(t *testing.T) {
	content := "# Title\n\nfirst paragraph\nstill first\n\nsecond paragraph mentions sourdough"
	g := graph.New()
	g.Add(1, "note", zb.Note{Content: content, Start: 2})
	book := &Notebook{Notes: g}

	got := book.Search("sourdough")

	if len(got) != 1 || got[0].Line != 7 {
		t.Fatalf("Expected a single match on line 7, got: %+v", got)
	}
}
