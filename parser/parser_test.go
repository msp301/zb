package parser

import (
	"reflect"
	"testing"
)

func TestParseValidNote(t *testing.T) {
	got := Parse("examples/valid_note.md")
	want := []Note{
		{
			Content: "This is a test file. #test\nWith a single note :)\n\n\\#hashtag\n\nLinks: [[202203311822]], [[202203310800]]",
			File:    "examples/valid_note.md",
			Id:      202001241300,
			Links:   []uint64{202203311822, 202203310800},
			Tags:    []string{"#test"},
			Title:   "This is a test file. #test",
		},
	}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("Unexpected note: %+v\n\nWanted: %+v", got, want)
	}
}

func TestParseCRLFNote(t *testing.T) {
	got := Parse("examples/windows_note.md")
	want := []Note{
		{
			Content: "Test Windows line endings.",
			File:    "examples/windows_note.md",
			Id:      202003092017,
			Links:   nil,
			Tags:    []string{"#Windows"},
			Title:   "Test Windows line endings.",
		},
	}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("Unexpected note: %+v\n\nWanted: %+v", got, want)
	}
}

func TestParseDuplicateTags(t *testing.T) {
	got := Parse("examples/tags.md")
	want := []Note{
		{
			Content: "This is a #test file. #test\nWith #several #tags mentioned\n#several times.\n\n## Subheading\n\n#Order is preserved.\nThis is a tag in quotes \"#1-1\"\n#C# and #C++ are valid tags.\nTag in a URL: http://example.com/#foo\n\n\\#ignored",
			File:    "examples/tags.md",
			Id:      202204192322,
			Links:   nil,
			Tags:    []string{"#test", "#several", "#tags", "#Order", "#C#", "#C++"},
			Title:   "Test tags mentioned multiple times",
		},
	}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("Unexpected note: %+v\n\nWanted: %+v", got, want)
	}
}

func TestParseMultiNoteFile(t *testing.T) {
	got := Parse("examples/multi_note.md")
	want := []Note{
		{
			Content: "talk about some #thing",
			File:    "examples/multi_note.md",
			Id:      202001300000,
			Links:   nil,
			Tags:    []string{"#thing"},
			Title:   "First note",
		},
		{
			Content: "#fruit\n* apple\n* orange",
			File:    "examples/multi_note.md",
			Id:      202002010000,
			Links:   []uint64{},
			Tags:    []string{"#fruit"},
			Title:   "Second note",
		},
		{
			Content: "No date or tags :(",
			File:    "examples/multi_note.md",
			Id:      0,
			Links:   []uint64{},
			Tags:    []string{},
			Title:   "Bad note",
		},
		{
			Content: "stuff --- stuff ---\n\n    A #link anchor should not be a tag\nhttp://localhost/test#anchor\n\n[ ] Don't forget #todo this important thing\n\n[[202002010000]]",
			File:    "examples/multi_note.md",
			Id:      202002020001,
			Links:   []uint64{202002010000},
			Tags:    []string{"#link", "#todo"},
			Title:   "Third note",
		},
		{
			Content: "Test note ID including time",
			File:    "examples/multi_note.md",
			Id:      202002272033,
			Links:   []uint64{},
			Tags:    []string{"#.Net"},
			Title:   "Test note ID including time",
		},
	}
	if len(got) != len(want) {
		t.Fatalf("Got %d notes but expected %d", len(got), len(want))
	}
	for i := range got {
		if !reflect.DeepEqual(got[i], want[i]) {
			t.Fatalf("Unexpected note %d: %+v\n\nWanted: %+v", i, got[i], want[i])
		}
	}
}
