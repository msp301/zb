package graph

import (
	"fmt"
	"reflect"
	"testing"
)

func TestNew(t *testing.T) {
	got := New()
	want := &Graph{
		Vertices:  map[uint64]Vertex{},
		Edges:     map[uint64]Edge{},
		Adjacency: map[uint64]map[uint64]int{},
	}

	if !reflect.DeepEqual(got, want) {
		t.Fatalf("Expected: %v\nGot: %v\n", want, got)
	}
}

func TestAddVertex(t *testing.T) {
	got := New()
	got.AddVertex(Vertex{Id: 1})
	got.AddVertex(Vertex{Id: 2})
	got.AddVertex(Vertex{Id: 3})

	want := &Graph{
		Vertices:  map[uint64]Vertex{1: {Id: 1}, 2: {Id: 2}, 3: {Id: 3}},
		Edges:     map[uint64]Edge{},
		Adjacency: map[uint64]map[uint64]int{},
	}

	if !reflect.DeepEqual(got, want) {
		t.Fatalf("Expected: %v\nGot: %v\n", want, got)
	}
}

func TestAddEdge(t *testing.T) {
	got := &Graph{
		Vertices:  map[uint64]Vertex{1: {Id: 1}, 2: {Id: 2}, 3: {Id: 3}},
		Edges:     map[uint64]Edge{},
		Adjacency: map[uint64]map[uint64]int{},
	}
	_ = got.AddEdge(Edge{Id: 1, From: 1, To: 2, Label: "link"})

	want := &Graph{
		Vertices:  map[uint64]Vertex{1: {Id: 1}, 2: {Id: 2}, 3: {Id: 3}},
		Edges:     map[uint64]Edge{1: {Id: 1, From: 1, To: 2, Label: "link"}, 2: {Id: 2, From: 2, To: 1, Label: "link"}},
		Adjacency: map[uint64]map[uint64]int{1: {2: 1}, 2: {1: 1}},
	}

	t.Logf("Edges: %v", got.Edges)

	if !reflect.DeepEqual(got, want) {
		t.Fatalf("Expected: %+v\nGot: %+v\n", want, got)
	}
}

func TestAddEdge_ErrorsOnNonExistentVertex(t *testing.T) {
	g := New()
	g.AddVertex(Vertex{Id: 1})
	err1 := g.AddEdge(Edge{Id: 101, From: 1, To: 2, Label: "link"})
	err2 := g.AddEdge(Edge{Id: 102, From: 2, To: 1, Label: "link"})

	want := fmt.Errorf("vertex does not exist: 2")

	if !reflect.DeepEqual(err1, want) {
		t.Fatalf("Error 1: Got '%v'", err1)
	}
	if !reflect.DeepEqual(err2, want) {
		t.Fatalf("Error 2: Got '%v'", err2)
	}
}

type isEdgeTest struct {
	A    uint64
	B    uint64
	Want bool
}

func TestIsEdge(t *testing.T) {
	graph := &Graph{
		Vertices:  map[uint64]Vertex{1: {}, 2: {}, 3: {}},
		Edges:     map[uint64]Edge{},
		Adjacency: map[uint64]map[uint64]int{1: {2: 1}, 2: {1: 1, 3: 2}, 3: {2: 1}},
	}

	tests := []isEdgeTest{
		{A: 1, B: 2, Want: true},
		{A: 2, B: 1, Want: true},
		{A: 3, B: 1, Want: false},
	}
	for _, test := range tests {
		got := graph.IsEdge(test.A, test.B)
		if got != test.Want {
			t.Fatalf("%v - Returned %v but expected %v", test, got, test.Want)
		}
	}
}

func TestWalk(t *testing.T) {
	graph := New()
	for i := 1; i <= 5; i++ {
		graph.AddVertex(Vertex{Id: uint64(i), Label: "vertex"})
	}
	_ = graph.AddEdge(Edge{From: 1, To: 2})
	_ = graph.AddEdge(Edge{From: 2, To: 5})
	_ = graph.AddEdge(Edge{From: 2, To: 3})

	var got []uint64
	graph.Walk(func(vertex Vertex, depth int) bool { got = append(got, vertex.Id); return true }, -1)

	want := []uint64{1, 2, 5, 3, 4}

	if !reflect.DeepEqual(got, want) {
		t.Fatalf("Expected: %v\nGot: %v\n", want, got)
	}
}
