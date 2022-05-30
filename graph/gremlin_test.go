package graph

import (
	"fmt"
	"reflect"
	"testing"
)

func TestTraversal(t *testing.T) {
	graph := New()
	got := Traversal(graph)
	want := &TraversalSource{
		graph:   graph,
		visited: map[uint64]bool{},
	}

	if !reflect.DeepEqual(got, want) {
		t.Fatalf("Expected: %v\nGot: %v\n", want, got)
	}
}

func TestV(t *testing.T) {
	traversal := Traversal(New()).V()
	if traversal.what != "vertex" {
		t.Fatalf("Got: %v\n", traversal.what)
	}
}

func TestIterate(t *testing.T) {
	g := New()
	g.AddVertex(Vertex{Id: 1, Label: "foo"})
	g.AddVertex(Vertex{Id: 2, Label: "bar"})
	g.AddVertex(Vertex{Id: 3, Label: "foo"})

	var vertices []uint64
	for vertex := range Traversal(g).V().Iterate() {
		vertices = append(vertices, vertex.Id)
	}

	want := []uint64{1, 2, 3}
	if !reflect.DeepEqual(vertices, want) {
		t.Fatalf("Expected: %v\nGot: %v\n", want, vertices)
	}
}

func TestNext(t *testing.T) {
	g := New()
	g.AddVertex(Vertex{Id: 1})
	g.AddVertex(Vertex{Id: 2})
	g.AddVertex(Vertex{Id: 3})

	var got []*Vertex
	traversal := Traversal(g).V()
	got = append(got, traversal.Next())
	got = append(got, traversal.Next())
	got = append(got, traversal.Next())
	got = append(got, traversal.Next())

	want := []*Vertex{{Id: 1}, {Id: 2}, {Id: 3}, nil}
	for i, vertex := range got {
		if vertex == nil && want[i] == nil {
			continue
		}
		if vertex == nil && want[i] != nil {
			t.Fatalf("Expected: %v\nGot: nil", want[i])
		}
		if !reflect.DeepEqual(vertex, want[i]) {
			t.Fatalf("Expected: %v\nGot: %v\n", want[i], vertex)
		}
	}
}

func TestHas(t *testing.T) {
	g := New()
	g.AddVertex(Vertex{Id: 1, Properties: map[string]interface{}{"Title": "foo", "Tags": "b"}})
	g.AddVertex(Vertex{Id: 2, Properties: map[string]interface{}{"Title": "bar", "Tags": []string{"a", "c"}}})
	g.AddVertex(Vertex{Id: 3, Properties: map[string]interface{}{"Title": "foo", "Tags": []string{"a", "b"}}})

	var vertices []uint64
	for vertex := range Traversal(g).V().Has("Tags", "c").Iterate() {
		vertices = append(vertices, vertex.Id)
	}

	want := []uint64{2}
	if !reflect.DeepEqual(vertices, want) {
		t.Fatalf("Expected: %v\nGot: %v\n", want, vertices)
	}

}

func TestHasLabel(t *testing.T) {
	g := New()
	g.AddVertex(Vertex{Id: 1, Label: "foo"})
	g.AddVertex(Vertex{Id: 2, Label: "bar"})
	g.AddVertex(Vertex{Id: 3, Label: "foo"})

	var vertices []uint64
	for vertex := range Traversal(g).V().HasLabel("foo").Iterate() {
		vertices = append(vertices, vertex.Id)
	}

	want := []uint64{1, 3}
	if !reflect.DeepEqual(vertices, want) {
		t.Fatalf("Expected: %v\nGot: %v\n", want, vertices)
	}
}

func TestValues(t *testing.T) {
	g := New()
	g.AddVertex(Vertex{Id: 1, Properties: map[string]interface{}{"value": "foo"}})
	g.AddVertex(Vertex{Id: 2})
	g.AddVertex(Vertex{Id: 3})

	var values []string
	for _, value := range Traversal(g).V().Values("value") {
		values = append(values, fmt.Sprint(value))
	}

	want := []string{"foo"}
	if !reflect.DeepEqual(values, want) {
		t.Fatalf("Expected: %v\nGot: %v\n", want, values)
	}
}
