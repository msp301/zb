package notebook

import (
	"github.com/msp301/graph"
	"reflect"
	"testing"
)

func TestSearchByTag(t *testing.T) {
	g := graph.New()
	g.Add(1, "note", nil)
	g.Add(2, "note", nil)
	g.Add(3, "tag", "#foo")
	g.Add(4, "note", nil)
	g.Add(5, "tag", "#bar")
	g.Edge(2, 3, "tag")
	g.Edge(1, 5, "tag")
	book := &Notebook{Notes: g}

	got := book.SearchByTags("#foo")
	want := []Result{{Line: -1, Value: graph.Vertex{Id: 2, Label: "note"}}}

	if !reflect.DeepEqual(got, want) {
		t.Fatalf("Expected: %+v\nGot: %+v", want, got)
	}
}
