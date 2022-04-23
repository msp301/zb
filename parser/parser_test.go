package parser

import (
	"reflect"
	"testing"
)

func TestParseValidNote(t *testing.T) {
	got := Parse("examples/valid_note.md")
	want := Note{
		Content: "This is a test file. #test\nWith a single note :)\n\\#hashtag\nLinks: [[202203311822]], [[202203310800]]",
		File:    "examples/valid_note.md",
		Id:      202001241300,
		Links:   []uint64{202203311822, 202203310800},
		Tags:    []string{"#test"},
		Title:   "This is a test file. #test",
	}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("Unexpected note: %+v", got)
	}
}
