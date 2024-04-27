package util

import (
	"regexp"
	"strings"
)

type ContextMatch struct {
	Text string
	Line int
}

var mdListRegex = regexp.MustCompile(`^(\s*)(?:\*|\+|-|\d+[.)])\s+`)

var mdListEntryRegex = regexp.MustCompile(`^(\s*)(?:(?:\*|\+|-|\d+[.)])\s+)?([^\n]+)`)

// cache contextRegex by input phrase
var contextRegexCache = make(map[string]*regexp.Regexp)

type ContextMatchFunc func(s string, phrase string) bool

func context(s string, phrase string, matchFunc ContextMatchFunc) ([]ContextMatch, bool) {
	contextRegex := contextRegex(phrase)
	matches := contextRegex.FindAllStringSubmatch(s, -1)
	if matches == nil {
		return nil, false
	}

	contexts := make([]ContextMatch, 0)

	for _, match := range matches {
		match := strings.TrimSpace(match[0])

		lineNumber := strings.Count(s[:strings.Index(s, match)], "\n") + 1

		if isMarkdownList(match) {
			for _, line := range strings.Split(match, "\n") {
				if matchFunc(line, phrase) {
					context := mdListEntryRegex.FindStringSubmatch(line)
					contexts = append(contexts, ContextMatch{Text: context[2], Line: lineNumber})
				}
			}
			continue
		}

		if isMarkdownTable(match) {
			for _, row := range strings.Split(match, "\n") {
				if matchFunc(row, phrase) {
					for _, cell := range strings.Split(row, "|") {
						if matchFunc(cell, phrase) {
							contexts = append(contexts, ContextMatch{Text: strings.TrimSpace(cell), Line: lineNumber})
						}
					}
				}
			}
			continue
		}

		contexts = append(contexts, ContextMatch{Text: match, Line: lineNumber})
	}

	return contexts, true
}

func isMarkdownList(line string) bool {
	return mdListRegex.MatchString(line)
}

func isMarkdownTable(line string) bool {
	return strings.HasPrefix(line, "|")
}

func contextRegex(phrase string) *regexp.Regexp {
	if contextRegexCache[phrase] == nil {
		input := regexp.QuoteMeta(phrase)
		contextRegexCache[phrase] = regexp.MustCompile(`(?i)(?:[^\n]\n?)*` + input + `(?:[^\n]\n?)*`)
	}

	return contextRegexCache[phrase]
}

func Context(s string, phrase string) ([]ContextMatch, bool) {
	return context(s, phrase, func(s string, t string) bool {
		return strings.Contains(s, phrase)
	})
}

func ContextFold(s string, phrase string) ([]ContextMatch, bool) {
	return context(s, phrase, func(s string, t string) bool {
		return strings.Contains(strings.ToLower(s), strings.ToLower(phrase))
	})
}
