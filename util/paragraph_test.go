package util

import (
	"fmt"
	"testing"
)

func TestParagraphMatches(t *testing.T) {
	tests := []struct {
		query   string
		content string
		want    bool
	}{
		{"string", "SomeString", true},
		{"string", "Some String", true},
		{"c", "C#", true},

		{"Foo", "|Column A|Column B|\n|------|------|\n|Value Foo|Value Bar|", true},
		{"Foo", "| Column A | Column B |\n| ------ | ------ |\n| Value Foo | Value Bar |", true},

		{"foo", "|Column A|Column B|\n|------|------|\n|Value Foo|Value Bar|", true},
		{"foo", "| Column A | Column B |\n| ------ | ------ |\n| Value Foo | Value Bar |", true},

		{"booom", "| Thing | Booom |\n", true},
	}

	for _, test := range tests {
		t.Run(fmt.Sprintf("ParagraphMatches('%v', '%v')", test.content, test.query), func(t *testing.T) {
			got := ParagraphMatches(test.content, test.query)
			if got != test.want {
				t.Fatalf("expected %t but was %t", test.want, got)
			}
		})
	}
}
