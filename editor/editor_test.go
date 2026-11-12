package editor

import (
	"reflect"
	"testing"
)

func TestArguments(t *testing.T) {
	tests := []struct {
		name     string
		template []string
		line     int
		want     []string
	}{
		{"no template appends the file", nil, 12, []string{"my note.md"}},
		{"flags without placeholders keep the file last", []string{"-w"}, 12, []string{"-w", "my note.md"}},
		{"line and file placeholders", []string{"+{line}", "{file}"}, 12, []string{"+12", "my note.md"}},
		{"placeholders within one argument", []string{"--goto", "{file}:{line}"}, 12, []string{"--goto", "my note.md:12"}},
		{"line placeholder alone still appends the file", []string{"+{line}"}, 12, []string{"+12", "my note.md"}},
		{"unknown line falls back to the first line", []string{"+{line}", "{file}"}, 0, []string{"+1", "my note.md"}},
	}

	for _, test := range tests {
		got := arguments(test.template, "my note.md", test.line)
		if !reflect.DeepEqual(got, test.want) {
			t.Errorf("%s: expected %v, got %v", test.name, test.want, got)
		}
	}
}
