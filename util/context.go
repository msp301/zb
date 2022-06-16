package util

import (
	"regexp"
	"strings"
)

func Context(s string, phrase string) string {
	input := regexp.QuoteMeta(phrase)
	contextRegex := regexp.MustCompile(`(?:[\w ]\n?)*` + input + `(?:[\w ]\n?)*`)
	matches := contextRegex.FindStringSubmatch(s)
	if len(matches) == 0 {
		return ""
	}

	return strings.TrimSpace(matches[0])
}
