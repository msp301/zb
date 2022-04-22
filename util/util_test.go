package util

import (
	"fmt"
	"testing"
)

func TestIsMetadata(t *testing.T) {
	tests := []struct {
		input []byte
		want  bool
	}{
		{[]byte("202004150100"), true},
		{[]byte("#this #that"), true},
		{[]byte(" #tag "), true},
		{[]byte("[[202203310100]]"), true},
		{[]byte("[[31/03/2022 01:00]]"), true},
		{[]byte("[[link to note]]"), false},
		{[]byte("15/04/2020 01:00"), false},
		{[]byte("# Heading 1"), false},
		{[]byte("## Heading 2"), false},
		{[]byte("Text with a #tag"), false},
		{[]byte("Just some text"), false},
	}

	for _, test := range tests {
		t.Run(fmt.Sprintf("IsMetadata('%s')", test.input), func(t *testing.T) {
			got := IsMetadata(test.input)
			if got != test.want {
				t.Fatalf("expected %t but was %t", test.want, got)
			}
		})

	}
}
