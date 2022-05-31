package util

import (
	"strings"
)

func Context(s string, phrase string) string {
	pos := strings.Index(s, phrase)
	if pos == -1 {
		return ""
	}

	start := pos
	end := pos + len(phrase) + 10

	return s[start:end]
}
