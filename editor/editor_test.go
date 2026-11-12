package editor

import (
	"reflect"
	"testing"
)

func TestLocationArgs(t *testing.T) {
	tests := []struct {
		program string
		line    int
		want    []string
	}{
		{"vim", 12, []string{"+12", "note.md"}},
		{"/usr/bin/nvim", 12, []string{"+12", "note.md"}},
		{"nano", 3, []string{"+3", "note.md"}},
		{"code", 12, []string{"--goto", "note.md:12"}},
		{"Code.cmd", 12, []string{"--goto", "note.md:12"}},
		{"hx", 12, []string{"note.md:12"}},
		{"notepad++.exe", 12, []string{"-n12", "note.md"}},
		{"notepad.exe", 12, []string{"note.md"}},
		{"unknown-editor", 12, []string{"note.md"}},
		{"vim", 0, []string{"note.md"}},
	}

	for _, test := range tests {
		got := locationArgs(test.program, "note.md", test.line)
		if !reflect.DeepEqual(got, test.want) {
			t.Errorf("%s line %d: expected %v, got %v", test.program, test.line, test.want, got)
		}
	}
}
