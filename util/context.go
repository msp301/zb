package util

import (
	"regexp"
	"strings"
)

func Context(s string, phrase string) ([]string, bool) {
	input := regexp.QuoteMeta(phrase)
	contextRegex := regexp.MustCompile(`(?i)(?:[^\n]\n?)*` + input + `(?:[^\n]\n?)*`)
	matches := contextRegex.FindAllStringSubmatch(s, -1)
	if matches == nil {
		return nil, false
	}

	contexts := make([]string, 0)

	for _, match := range matches {
		match := strings.TrimSpace(match[0])
		if isMarkdownList(match) {
			for _, line := range strings.Split(match, "\n") {
				if strings.Contains(line, phrase) {
					mdListRegex := regexp.MustCompile(`^(\s*)(?:(?:\*|\+|-|\d+[.)])\s+)?([^\n]+)`)
					context := mdListRegex.FindStringSubmatch(line)
					contexts = append(contexts, context[2])
				}
			}
			continue
		}
		contexts = append(contexts, match)
	}

	return contexts, true
}

func isMarkdownList(line string) bool {
	mdListRegex := regexp.MustCompile(`^(\s*)(?:\*|\+|-|\d+[.)])\s+`)
	return mdListRegex.MatchString(line)
}
