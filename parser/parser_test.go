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
			Start:   3,
			Tags:    []string{"#test"},
			Title:   "This is a test file. #test",
		},
	}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("Unexpected note: %+v\n\nWanted: %+v", got, want)
	}
}

func TestParseWikilinks(t *testing.T) {
	got := Parse("examples/wikilinks.md")
	want := []Note{
		{
			Content: "# Test wikilinks\n\nTest note with different [[202405011240|link styles]]. Another link ([[202503021105]]\n\nand [[202203310800]] but not [[foobar]] or [202501012210]",
			File:    "examples/wikilinks.md",
			Id:      202506162200,
			Links:   []uint64{202405011240, 202503021105, 202203310800},
			Start:   3,
			Tags:    []string{},
			Title:   "Test wikilinks",
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
			Start:   4,
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
			Content: "# Test tags mentioned multiple times\n\nThis is a #test file. #test\nWith #several #tags mentioned\n#several times.\n\n## Subheading\n\n#Order is preserved.\nThis is a tag in quotes \"#1-1\"\n#C# and #C++ are valid tags.\nTag in a URL: http://example.com/#foo\n\n\\#ignored",
			File:    "examples/tags.md",
			Id:      202204192322,
			Links:   nil,
			Start:   3,
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
			Content: "# First note\n\ntalk about some #thing",
			File:    "examples/multi_note.md",
			Id:      202001300000,
			Links:   nil,
			Start:   4,
			Tags:    []string{"#thing"},
			Title:   "First note",
		},
		{
			Content: "# Second note\n\n#fruit\n* apple\n* orange",
			File:    "examples/multi_note.md",
			Id:      202002010000,
			Links:   []uint64{},
			Start:   12,
			Tags:    []string{"#fruit"},
			Title:   "Second note",
		},
		{
			Content: "# Bad note\n\nNo date or tags :(",
			File:    "examples/multi_note.md",
			Id:      0,
			Links:   []uint64{},
			Start:   20,
			Tags:    []string{},
			Title:   "Bad note",
		},
		{
			Content: "# Third note\n\nstuff --- stuff ---\n\n    A #link anchor should not be a tag\nhttp://localhost/test#anchor\n\n[ ] Don't forget #todo this important thing\n\n[[202002010000]]",
			File:    "examples/multi_note.md",
			Id:      202002020001,
			Links:   []uint64{202002010000},
			Start:   28,
			Tags:    []string{"#link", "#todo"},
			Title:   "Third note",
		},
		{
			Content: "Test note ID including time",
			File:    "examples/multi_note.md",
			Id:      202002272033,
			Links:   []uint64{},
			Start:   44,
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
