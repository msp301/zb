package util

import (
	"reflect"
	"testing"
)

func TestCommandSplitsArguments(t *testing.T) {
	got := Command("zb-missing-program -F -R", "note.md").Args
	want := []string{"zb-missing-program", "-F", "-R", "note.md"}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("expected %v, got %v", want, got)
	}
}

func TestCommandWithoutArguments(t *testing.T) {
	got := Command("zb-missing-program", "note.md").Args
	want := []string{"zb-missing-program", "note.md"}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("expected %v, got %v", want, got)
	}
}
