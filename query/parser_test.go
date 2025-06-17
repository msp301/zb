package query

import (
	"testing"

	"github.com/msp301/graph"
)

func TestParse(t *testing.T) {
	want := &graph.Graph{
		Type: graph.DIRECTED,
		Vertices: map[uint64]graph.Vertex{
			1: *tokenVertex(1, TERM, "test"),
			2: *tokenVertex(2, TERM, "a"),
			3: *tokenVertex(3, TERM, "thing"),
		},
		Edges:     map[uint64]graph.Edge{},
		Adjacency: map[uint64]map[uint64]int{},
	}

	check(t, "test a thing", want)
}

func check(t *testing.T, input string, want *graph.Graph) {
	lexer := New(input)
	parser := NewParser(lexer)

	got := parser.Parse()

	if len(got.Vertices) != len(want.Vertices) {
		t.Fatalf("Expected: %d vertices. Got: %d verticies\n", len(want.Vertices), len(got.Vertices))
	}

	for key, gotVertex := range got.Vertices {
		wantVertex := want.Vertices[key]
		if gotVertex.Id != wantVertex.Id {
			t.Fatalf("Expected %d.Id: %v. Got: %v\n", key, wantVertex.Id, gotVertex.Id)
		}

		if gotVertex.Label != wantVertex.Label {
			t.Fatalf("Expected %d.Label: %v. Got: %v\n", key, wantVertex.Label, gotVertex.Label)
		}

		switch gotVal := gotVertex.Value.(type) {
		case Token:
			wantVal := want.Vertices[key].Value.(*Token)

			if gotVal.Type != wantVal.Type {
				t.Fatalf("Expected %d.Value.Type: %v. Got: %v\n", key, wantVal.Type, gotVal.Type)
			}

			if gotVal.Value != wantVal.Value {
				t.Fatalf("Expected %d.Value: %v. Got: %v\n", key, wantVal.Value, gotVal.Value)
			}
		default:
			t.Fatalf("Unexpected Token %d: Got: %v\n", key, got)
		}
	}
}

func tokenVertex(id uint64, tokenType TokenType, value string) *graph.Vertex {
	return &graph.Vertex{
		Id:    id,
		Label: tokenString(tokenType),
		Value: &Token{Type: tokenType, Value: value},
	}
}
