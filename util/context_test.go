package util

import (
	"fmt"
	"reflect"
	"testing"
)

func TestContext(t *testing.T) {

	tests := []struct {
		source string
		phrase string
		want   []string
	}{
		{"This is some example text about nothing or maybe something", "text", []string{"This is some example text about nothing or maybe something"}},
		{"This is some example text about nothing\nor maybe something", "text", []string{"This is some example text about nothing\nor maybe something"}},
		{"This is some example text about nothing\n\nor maybe something", "text", []string{"This is some example text about nothing"}},

		{"1. This is an example list about nothing\nor maybe something", "list", []string{"This is an example list about nothing"}},
		{"1) This is an example list about nothing\nor maybe something", "list", []string{"This is an example list about nothing"}},
		{"+ This is an example list about nothing\nor maybe something", "list", []string{"This is an example list about nothing"}},
		{"- This is an example list about nothing\nor maybe something", "list", []string{"This is an example list about nothing"}},
		{"* This is an example list about nothing\nor maybe something", "list", []string{"This is an example list about nothing"}},
		{"   * This is an example list about nothing\nor maybe something", "list", []string{"This is an example list about nothing"}},
		{" * This is an example\n* list about nothing\nor maybe something", "list", []string{"list about nothing"}},

		{"Example 1\n\nanother example\n\nand another", "another", []string{"another example", "and another"}},
	}

	for _, test := range tests {
		t.Run(fmt.Sprintf("Context('%v', '%v')", test.source, test.phrase), func(t *testing.T) {
			got, ok := Context(test.source, test.phrase)
			if !ok {
				t.Fatalf("Expected ok but was not ok")
			}
			if !reflect.DeepEqual(got, test.want) {
				t.Fatalf("expected '%s' but was '%s'", test.want, got)
			}
		})

	}
}
