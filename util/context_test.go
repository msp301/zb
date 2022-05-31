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
