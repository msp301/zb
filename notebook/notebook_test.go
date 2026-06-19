package notebook

import (
	"reflect"
	"testing"

	"github.com/msp301/graph"
)

func TestMatchedTags(t *testing.T) {
	g := graph.New()
	g.Add(1, "note", nil)
	g.Add(2, "note", nil)
	g.Add(3, "note", nil)

	g.Add(4, "tag", "#test")
	g.Add(5, "tag", "#LLM")
	g.Add(6, "tag", "#LargeLanguageModel")
	g.Add(7, "tag", "#LargeLanguageModels")

	book := &Notebook{Notes: g}

	got := book.MatchedTags("llm")
	want := []matchedTag{
		{Term: "llm", Distance: 1, Tag: "#LLM", Vertex: graph.Vertex{Id: 5, Label: "tag", Value: "#LLM"}},
		{Term: "llm", Distance: 16, Tag: "#LargeLanguageModel", Vertex: graph.Vertex{Id: 6, Label: "tag", Value: "#LargeLanguageModel"}},
		{Term: "llm", Distance: 17, Tag: "#LargeLanguageModels", Vertex: graph.Vertex{Id: 7, Label: "tag", Value: "#LargeLanguageModels"}},
	}

	if !reflect.DeepEqual(got, want) {
		t.Fatalf("Expected: %+v\nGot: %+v", want, got)
	}
}

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

func TestSearchByTags_PrioritisesMatches(t *testing.T) {
	g := graph.New()
	g.Add(1, "note", nil)
	g.Add(2, "note", nil)
	g.Add(3, "note", nil)

	g.Add(4, "tag", "#test")
	g.Add(5, "tag", "#LLM")
	g.Add(6, "tag", "#LargeLanguageModel")
	g.Add(7, "tag", "#LargeLanguageModels")

	g.Edge(2, 4, "tag")
	g.Edge(3, 4, "tag")
	g.Edge(3, 7, "tag")

	book := &Notebook{Notes: g}

	got := book.SearchByTags("test", "llm")
	want := []Result{{Line: -1, Value: graph.Vertex{Id: 3, Label: "note"}}}

	if !reflect.DeepEqual(got, want) {
		t.Fatalf("Expected: %+v\nGot: %+v", want, got)
	}
}

func TestTagIntersection(t *testing.T) {
	g := graph.New()
	g.Add(1, "note", nil)
	g.Add(2, "note", nil)
	g.Add(3, "note", nil)

	g.Add(4, "tag", "#test")
	g.Add(5, "tag", "#LLM")
	g.Add(6, "tag", "#LargeLanguageModel")
	g.Add(7, "tag", "#LargeLanguageModels")

	g.Edge(2, 4, "tag")
	g.Edge(3, 4, "tag")
	g.Edge(3, 7, "tag")

	book := &Notebook{Notes: g}
	tags := []matchedTag{
		{Term: "test", Distance: 1, Tag: "#test", Vertex: graph.Vertex{Id: 4, Label: "tag", Value: "#test"}},
		{Term: "llm", Distance: 1, Tag: "#LLM", Vertex: graph.Vertex{Id: 5, Label: "tag", Value: "#LLM"}},
		{Term: "llm", Distance: 16, Tag: "#LargeLanguageModel", Vertex: graph.Vertex{Id: 6, Label: "tag", Value: "#LargeLanguageModel"}},
		{Term: "llm", Distance: 17, Tag: "#LargeLanguageModels", Vertex: graph.Vertex{Id: 7, Label: "tag", Value: "#LargeLanguageModels"}},
	}

	got := book.TagIntersection(tags)
	want := []graph.Vertex{{Id: 3, Label: "note"}}

	if !reflect.DeepEqual(got, want) {
		t.Fatalf("Expected: %+v\nGot: %+v", want, got)
	}
}

func TestTagIntersection_MatchesNoteWithSecondClosestTagForTerm(t *testing.T) {
	g := graph.New()
	g.Add(1, "note", nil)
	g.Add(2, "note", nil)

	g.Add(3, "tag", "#foo")
	g.Add(4, "tag", "#bar")
	g.Add(5, "tag", "#foobar")

	g.Edge(1, 3, "tag")
	g.Edge(1, 5, "tag")
	g.Edge(2, 4, "tag")

	book := &Notebook{Notes: g}
	tags := []matchedTag{
		{Term: "foo", Distance: 1, Tag: "#foo", Vertex: graph.Vertex{Id: 3, Label: "tag", Value: "#foo"}},
		{Term: "bar", Distance: 1, Tag: "#bar", Vertex: graph.Vertex{Id: 4, Label: "tag", Value: "#bar"}},
		{Term: "bar", Distance: 4, Tag: "#foobar", Vertex: graph.Vertex{Id: 5, Label: "tag", Value: "#foobar"}},
	}

	got := book.TagIntersection(tags)
	want := []graph.Vertex{{Id: 1, Label: "note"}}

	if !reflect.DeepEqual(got, want) {
		t.Fatalf("Expected: %+v\nGot: %+v", want, got)
	}
}

func TestTagIntersection_DoesNotCountSameSearchTermTwice(t *testing.T) {
	g := graph.New()
	g.Add(1, "note", nil) // has #golang + #gopher (both match "go", no "web")
	g.Add(2, "note", nil) // has #golang + #web (correct match)

	g.Add(3, "tag", "#golang")
	g.Add(4, "tag", "#gopher")
	g.Add(5, "tag", "#web")

	g.Edge(1, 3, "tag")
	g.Edge(1, 4, "tag")
	g.Edge(2, 3, "tag")
	g.Edge(2, 5, "tag")

	book := &Notebook{Notes: g}
	tags := []matchedTag{
		{Term: "go", Distance: 1, Tag: "#golang", Vertex: graph.Vertex{Id: 3, Label: "tag", Value: "#golang"}},
		{Term: "go", Distance: 2, Tag: "#gopher", Vertex: graph.Vertex{Id: 4, Label: "tag", Value: "#gopher"}},
		{Term: "web", Distance: 1, Tag: "#web", Vertex: graph.Vertex{Id: 5, Label: "tag", Value: "#web"}},
	}

	got := book.TagIntersection(tags)
	want := []graph.Vertex{{Id: 2, Label: "note"}}

	if !reflect.DeepEqual(got, want) {
		t.Fatalf("Expected: %+v\nGot: %+v", want, got)
	}
}
