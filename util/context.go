package util

import (
	"regexp"
	"strings"
)

func Context(s string, phrase string) string {
	input := regexp.QuoteMeta(phrase)
	contextRegex := regexp.MustCompile(`(?:[^\n]\n?)*` + input + `(?:[^\n]\n?)*`)
	matches := contextRegex.FindStringSubmatch(s)
	if len(matches) == 0 {
		return ""
	}

	match := strings.TrimSpace(matches[0])
	if isMarkdownList(match) {
		for _, line := range strings.Split(match, "\n") {
			if strings.Contains(line, phrase) {
				mdListRegex := regexp.MustCompile(`^(\s*)(?:\*|\+|-|\d+[.)])\s+([^\n]+)`)
				context := mdListRegex.FindStringSubmatch(line)
				return context[2]
			}
		}
	}

	return match
}

func isMarkdownList(line string) bool {
	mdListRegex := regexp.MustCompile(`^(\s*)(?:\*|\+|-|\d+[.)])\s+`)
	return mdListRegex.MatchString(line)
}
