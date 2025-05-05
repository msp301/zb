package notebook

import (
	"github.com/msp301/graph"
	"reflect"
	"testing"
)

func TestSearchByTag(t *testing.T) {
	g := graph.New()
	g.AddVertex(graph.Vertex{Id: 1, Label: "note"})
	g.AddVertex(graph.Vertex{Id: 2, Label: "note"})
	g.AddVertex(graph.Vertex{Id: 3, Label: "tag", Properties: map[string]interface{}{"Value": "#foo"}})
	g.AddVertex(graph.Vertex{Id: 4, Label: "note"})
	g.AddVertex(graph.Vertex{Id: 5, Label: "tag", Properties: map[string]interface{}{"Value": "#bar"}})
	g.AddEdge(graph.Edge{Id: 100, Label: "tag", From: 2, To: 3})
	g.AddEdge(graph.Edge{Id: 101, Label: "tag", From: 1, To: 5})
	book := &Notebook{Notes: g}

	got := book.SearchByTags("#foo")
	want := []Result{{Line: -1, Value: graph.Vertex{Id: 2, Label: "note"}}}

	if !reflect.DeepEqual(got, want) {
		t.Fatalf("Expected: %+v\nGot: %+v", want, got)
	}
}
