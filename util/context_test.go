package util

import (
	"fmt"
	"testing"
)

func TestContext(t *testing.T) {

	tests := []struct {
		source string
		phrase string
		want   string
	}{
		{"This is some example text about nothing or maybe something", "text", "This is some example text about nothing or maybe something"},
		{"This is some example text about nothing\nor maybe something", "text", "This is some example text about nothing\nor maybe something"},
		{"This is some example text about nothing\n\nor maybe something", "text", "This is some example text about nothing"},

		{"1. This is an example list about nothing\nor maybe something", "list", "This is an example list about nothing"},
		{"1) This is an example list about nothing\nor maybe something", "list", "This is an example list about nothing"},
		{"+ This is an example list about nothing\nor maybe something", "list", "This is an example list about nothing"},
		{"- This is an example list about nothing\nor maybe something", "list", "This is an example list about nothing"},
		{"* This is an example list about nothing\nor maybe something", "list", "This is an example list about nothing"},
		{"   * This is an example list about nothing\nor maybe something", "list", "This is an example list about nothing"},
		{" * This is an example\n* list about nothing\nor maybe something", "list", "list about nothing"},
	}

	for _, test := range tests {
		t.Run(fmt.Sprintf("Context('%v', '%v')", test.source, test.phrase), func(t *testing.T) {
			got := Context(test.source, test.phrase)
			if got != test.want {
				t.Fatalf("expected '%s' but was '%s'", test.want, got)
			}
		})

	}
}
