package notebook

import (
	"reflect"
	"testing"

	"github.com/msp301/graph"
)

func bookWhereOnlyOneNoteHasEveryTag() *Notebook {
	g := graph.New()
	g.Add(1, "note", nil)
	g.Add(2, "note", nil)
	g.Add(3, "note", nil)
	g.Add(10, "tag", "#foo")
	g.Add(11, "tag", "#bar")
	g.Add(12, "tag", "#other")

	tagNote(g, 1, 10, 11)
	tagNote(g, 2, 10, 12)
	tagNote(g, 3, 12)

	return &Notebook{Notes: g}
}

func TestTagQuerySubgraph_IncludesNotesOfMatchedTagsAndAllTheirTags(t *testing.T) {
	subgraph := bookWhereOnlyOneNoteHasEveryTag().tagQuerySubgraph("foo", "bar")

	gotVertices := subgraph.sortedVertexIds()
	wantVertices := []uint64{1, 2, 10, 11, 12}
	if !reflect.DeepEqual(gotVertices, wantVertices) {
		t.Fatalf("Expected vertices: %v\nGot: %v", wantVertices, gotVertices)
	}

	gotEdges := subgraph.sortedEdges()
	wantEdges := []tagQueryEdge{{1, 10}, {1, 11}, {2, 10}, {2, 12}}
	if !reflect.DeepEqual(gotEdges, wantEdges) {
		t.Fatalf("Expected edges: %v\nGot: %v", wantEdges, gotEdges)
	}
}

func TestTagQuerySubgraph_MarksOnlyNotesReturnedBySearch(t *testing.T) {
	subgraph := bookWhereOnlyOneNoteHasEveryTag().tagQuerySubgraph("foo", "bar")

	want := map[uint64]struct{}{1: {}}
	if !reflect.DeepEqual(subgraph.returned, want) {
		t.Fatalf("Expected returned notes: %v\nGot: %v", want, subgraph.returned)
	}
}

func TestTagQuerySubgraph_MarksNothingReturnedWhenATermHasNoMatchingTag(t *testing.T) {
	subgraph := bookWhereOnlyOneNoteHasEveryTag().tagQuerySubgraph("foo", "missing")

	if len(subgraph.returned) != 0 {
		t.Fatalf("Expected no returned notes, got: %v", subgraph.returned)
	}
	if _, ok := subgraph.vertices[10]; !ok {
		t.Fatalf("Expected matched tag to still be drawn")
	}
}

func TestVisualizeTagQuery(t *testing.T) {
	got := bookWhereOnlyOneNoteHasEveryTag().VisualizeTagQuery("foo", "bar")
	want := `digraph {
  1 [label="note:1" shape=box style=filled fillcolor=lightgreen];
  2 [label="note:2" shape=box];
  10 [label="#foo\nterm: foo, distance: 1" shape=ellipse style=filled fillcolor=lightblue];
  11 [label="#bar\nterm: bar, distance: 1" shape=ellipse style=filled fillcolor=lightblue];
  12 [label="#other" shape=ellipse];
  1 -> 10;
  1 -> 11;
  2 -> 10;
  2 -> 12;
}
`
	if got != want {
		t.Fatalf("Expected:\n%s\nGot:\n%s", want, got)
	}
}

func TestDotString_EscapesQuotesAndBackslashes(t *testing.T) {
	got := dotString(`a "quoted" \ tag`, "second")
	want := `"a \"quoted\" \\ tag\nsecond"`
	if got != want {
		t.Fatalf("Expected: %s\nGot: %s", want, got)
	}
}
