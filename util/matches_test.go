package util

import (
	"fmt"
	"testing"
)

func TestMatches(t *testing.T) {

	tests := []struct {
		search string
		target string
		want   bool
	}{
		{"string", "SomeString", true},
		{"String", "SomeString", true},
		{"t", "this starts with character", true},
		{"t", "doesn't start with character", false},
		{"s", "doesn't start with character", false},
		{"st", "doesn't start with character", false},
		{"sta", "doesn't start with character", true},
		{"else", "SomeString", false},
		{"c", "#C#", true},
	}

	for _, test := range tests {
		t.Run(fmt.Sprintf("Matches('%v', '%v')", test.search, test.target), func(t *testing.T) {
			got := Matches(test.search, test.target)
			if got != test.want {
				t.Fatalf("expected %t but was %t", test.want, got)
			}
		})

	}
}
