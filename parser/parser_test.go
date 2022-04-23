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
		t.Fatalf("Unexpected note: %+v\n\nWanted: %+v", got, want)
	}
}

func TestParseCRLFNote(t *testing.T) {
	got := Parse("examples/windows_note.md")
	want := Note{
		Content: "Test Windows line endings.",
		File:    "examples/windows_note.md",
		Id:      202003092017,
		Links:   []uint64{},
		Tags:    []string{"#Windows"},
		Title:   "Test Windows line endings.",
	}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("Unexpected note: %+v\n\nWanted: %+v", got, want)
	}
}

func TestParseDuplicateTags(t *testing.T) {
	got := Parse("examples/tags.md")
	want := Note{
		Content: "This is a #test file. #test\nWith #several #tags mentioned\n#several times.\n#Order is preserved.\n\\#ignored",
		File:    "examples/tags.md",
		Id:      202204192322,
		Links:   []uint64{},
		Tags:    []string{"#test", "#several", "#tags", "#Order"},
		Title:   "Test tags mentioned multiple times",
	}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("Unexpected note: %+v\n\nWanted: %+v", got, want)
	}
}
