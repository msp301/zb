package query

import (
	"fmt"
	"testing"
)

func TestNextToken(t *testing.T) {
	tests := []Token{
		{Type: TERM, Value: "bar"},
		{Type: LEFT_BRACKET, Value: "("},
		{Type: TERM, Value: "foo"},
		{Type: OR, Value: "or"},
		{Type: TERM, Value: "baz"},
		{Type: RIGHT_BRACKET, Value: ")"},
		{Type: AND, Value: "and"},
		{Type: NOT, Value: "not"},
		{Type: DOUBLE_QUOTE, Value: "\""},
		{Type: TERM, Value: "bop"},
		{Type: DOUBLE_QUOTE, Value: "\""},
		{Type: OR, Value: "or"},
		{Type: SINGLE_QUOTE, Value: "'"},
		{Type: TERM, Value: "pop"},
		{Type: SINGLE_QUOTE, Value: "'"},
		{Type: END, Value: ""},
	}

	lexer := New("bar (foo or baz) and not \"bop\" or 'pop'")

	for index, want := range tests {
		t.Run(fmt.Sprintf("Position %d", index), func(t *testing.T) {
			got := lexer.NextToken()
			if got.Type != want.Type {
				t.Fatalf("'%s' expected type '%v' but was '%v'", got.Value, want.Type, got.Type)
			}
			if got.Value != want.Value {
				t.Fatalf("expected value '%s' but was '%s'", want.Value, got.Value)
			}
		})
	}
}
