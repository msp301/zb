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
		want   []ContextMatch
	}{
		{"This is some example text about nothing or maybe something", "text", []ContextMatch{{Text: "This is some example text about nothing or maybe something", Line: 1}}},
		{"This is some example text about nothing\nor maybe something", "text", []ContextMatch{{Text: "This is some example text about nothing\nor maybe something", Line: 1}}},
		{"This is some example text about nothing\n\nor maybe something", "text", []ContextMatch{{Text: "This is some example text about nothing", Line: 1}}},

		{"1. This is an example list about nothing\nor maybe something", "list", []ContextMatch{{Text: "This is an example list about nothing", Line: 1}}},
		{"1) This is an example list about nothing\nor maybe something", "list", []ContextMatch{{Text: "This is an example list about nothing", Line: 1}}},
		{"+ This is an example list about nothing\nor maybe something", "list", []ContextMatch{{Text: "This is an example list about nothing", Line: 1}}},
		{"- This is an example list about nothing\nor maybe something", "list", []ContextMatch{{Text: "This is an example list about nothing", Line: 1}}},
		{"* This is an example list about nothing\nor maybe something", "list", []ContextMatch{{Text: "This is an example list about nothing", Line: 1}}},
		{"   * This is an example list about nothing\nor maybe something", "list", []ContextMatch{{Text: "This is an example list about nothing", Line: 1}}},
		{" * This is an example\n* list about nothing\nor maybe something", "list", []ContextMatch{{Text: "list about nothing", Line: 1}}},

		{" * This is an example\n* list about nothing\nNot a list", "list", []ContextMatch{{Text: "list about nothing", Line: 1}, {Text: "Not a list", Line: 2}}},

		{" * This is a list entry\n* list about nothing\nor maybe something", "list", []ContextMatch{{Text: "This is a list entry", Line: 1}, {Text: "list about nothing", Line: 2}}},
		{"Example 1\n\nanother example\n\nand another", "another", []ContextMatch{{Text: "another example", Line: 1}, {Text: "and another", Line: 3}}},

		{"|Column A|Column B|\n|------|------|\n|Value foo|Value bar|", "foo", []ContextMatch{{Text: "Value foo", Line: 1}}},
		{"| Column A | Column B |\n| ------ | ------ |\n| Value foo | Value bar |", "foo", []ContextMatch{{Text: "Value foo", Line: 1}}},
	}

	for _, test := range tests {
		t.Run(fmt.Sprintf("Context('%v', '%v')", test.source, test.phrase), func(t *testing.T) {
			got, ok := Context(test.source, test.phrase)
			if !ok {
				t.Fatalf("Expected ok but was not ok")
			}
			if !reflect.DeepEqual(got, test.want) {
				t.Fatalf("expected '%+v' but was '%+v'", test.want, got)
			}
		})
	}
}

func TestContextFold(t *testing.T) {

	tests := []struct {
		source string
		phrase string
		want   []ContextMatch
	}{
		{"|Column A|Column B|\n|------|------|\n|Value foo|Value bar|", "foo", []ContextMatch{{Text: "Value foo", Line: 3}}},
		{"| Column A | Column B |\n| ------ | ------ |\n| Value foo | Value bar |", "foo", []ContextMatch{{Text: "Value foo", Line: 3}}},

		{"|Column A|Column B|\n|------|------|\n|Value Foo|Value Bar|", "foo", []ContextMatch{{Text: "Value Foo", Line: 3}}},
		{"| Column A | Column B |\n| ------ | ------ |\n| Value Foo | Value Bar |", "foo", []ContextMatch{{Text: "Value Foo", Line: 3}}},
	}

	for _, test := range tests {
		t.Run(fmt.Sprintf("ContextFold('%v', '%v')", test.source, test.phrase), func(t *testing.T) {
			got, ok := ContextFold(test.source, test.phrase)
			if !ok {
				t.Fatalf("Expected ok but was not ok")
			}
			if !reflect.DeepEqual(got, test.want) {
				t.Fatalf("expected '%+v' but was '%+v'", test.want, got)
			}
		})
	}
}
