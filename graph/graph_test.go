package graph

import (
	"reflect"
	"testing"
)

func TestNew(t *testing.T) {
	got := New()
	want := &Graph{
		Vertices: map[uint64]bool{},
		Edges:    map[uint64][]uint64{},
	}

	if !reflect.DeepEqual(got, want) {
		t.Fatalf("Expected: %v\nGot: %v\n", want, got)
	}
}

func TestAddVertex(t *testing.T) {
	got := New()
	got.AddVertex(1)
	got.AddVertex(2)
	got.AddVertex(3)

	want := &Graph{
		Vertices: map[uint64]bool{1: true, 2: true, 3: true},
		Edges:    map[uint64][]uint64{},
	}

	if !reflect.DeepEqual(got, want) {
		t.Fatalf("Expected: %v\nGot: %v\n", want, got)
	}
}

func TestAddEdge(t *testing.T) {
	got := &Graph{
		Vertices: map[uint64]bool{1: true, 2: true, 3: true},
		Edges:    map[uint64][]uint64{},
	}
	got.AddEdge(1, 2)

	want := &Graph{
		Vertices: map[uint64]bool{1: true, 2: true, 3: true},
		Edges:    map[uint64][]uint64{1: {2}, 2: {1}},
	}

	if !reflect.DeepEqual(got, want) {
		t.Fatalf("Expected: %v\nGot: %v\n", want, got)
	}
}

type isEdgeTest struct {
	A    uint64
	B    uint64
	Want bool
}

func TestIsEdge(t *testing.T) {
	graph := &Graph{
		Vertices: map[uint64]bool{1: true, 2: true, 3: true},
		Edges:    map[uint64][]uint64{1: {2}, 2: {1, 3}, 3: {2}},
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
