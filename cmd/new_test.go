package cmd

import (
	"os"
	"testing"
)

func TestCreateNoteDoesNotOverwriteExistingNote(t *testing.T) {
	dir := t.TempDir()

	first, err := createNote(dir)
	if err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(first, []byte("existing content"), 0666); err != nil {
		t.Fatal(err)
	}

	second, err := createNote(dir)
	if err != nil {
		t.Fatal(err)
	}
	if second == first {
		t.Fatalf("expected a new note path, got %s twice", first)
	}

	content, err := os.ReadFile(first)
	if err != nil {
		t.Fatal(err)
	}
	if string(content) != "existing content" {
		t.Fatalf("expected first note to be untouched, got %q", content)
	}
}
