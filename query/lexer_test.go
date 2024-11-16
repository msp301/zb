package query

import (
	"fmt"
	"testing"
)

func TestNextToken(t *testing.T) {
	tests := []Token{
		{Type: TERM, Value: "bar"},
	}

	lexer := New("bar")

	for index, want := range tests {
		t.Run(fmt.Sprintf("Position %d", index), func(t *testing.T) {
			got := lexer.NextToken()
			if got.Type != want.Type {
				t.Fatalf("expected type %v but was %v", want.Type, got.Type)
			}
			if got.Value != want.Value {
				t.Fatalf("expected value %s but was %s", want.Value, got.Value)
			}
		})
	}
}
