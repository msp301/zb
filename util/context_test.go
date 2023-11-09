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

		{" * This is an example\n* list about nothing\nNot a list", "list", []string{"list about nothing", "Not a list"}},

		{" * This is a list entry\n* list about nothing\nor maybe something", "list", []string{"This is a list entry", "list about nothing"}},
		{"Example 1\n\nanother example\n\nand another", "another", []string{"another example", "and another"}},

		{"|Column A|Column B|\n|------|------|\n|Value foo|Value bar|", "foo", []string{"Value foo"}},
		{"| Column A | Column B |\n| ------ | ------ |\n| Value foo | Value bar |", "foo", []string{"Value foo"}},
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

func TestContextFold(t *testing.T) {

	tests := []struct {
		source string
		phrase string
		want   []string
	}{
		{"|Column A|Column B|\n|------|------|\n|Value foo|Value bar|", "foo", []string{"Value foo"}},
		{"| Column A | Column B |\n| ------ | ------ |\n| Value foo | Value bar |", "foo", []string{"Value foo"}},

		{"|Column A|Column B|\n|------|------|\n|Value Foo|Value Bar|", "foo", []string{"Value Foo"}},
		{"| Column A | Column B |\n| ------ | ------ |\n| Value Foo | Value Bar |", "foo", []string{"Value Foo"}},
	}

	for _, test := range tests {
		t.Run(fmt.Sprintf("ContextFold('%v', '%v')", test.source, test.phrase), func(t *testing.T) {
			got, ok := ContextFold(test.source, test.phrase)
			if !ok {
				t.Fatalf("Expected ok but was not ok")
			}
			if !reflect.DeepEqual(got, test.want) {
				t.Fatalf("expected '%s' but was '%s'", test.want, got)
			}
		})
	}
}
