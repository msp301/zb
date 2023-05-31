package notebook

import (
	"github.com/lithammer/fuzzysearch/fuzzy"
	"github.com/msp301/zb/graph"
	"github.com/msp301/zb/parser"
	"regexp"
	"strings"
)

func (book *Notebook) Search(query string) []Result {
	results := []Result{}
	traversal := graph.Traversal(book.Notes)
	for vertex := range traversal.V().HasLabel("note").Iterate() {
		var context []string
		matched := false
		switch val := vertex.Properties["Value"].(type) {
		case parser.Note:
			paragraphs := extractParagraphs(val.Content)
			for _, paragraph := range paragraphs {
				if matches(paragraph, query) {
					context = append(context, paragraph)
					matched = true
					break
				}
			}

			if !matched {
				continue
			}
		}

		for _, context := range context {
			result := Result{
				Context: context,
				Value:   vertex,
			}
			results = append(results, result)
		}
	}
	return results
}

func extractParagraphs(content string) []string {
	return strings.Split(content, "\n\n")
}

func matches(content string, query string) bool {
	tokens := strings.Fields(content)
	for _, token := range tokens {
		if len(query) > 3 && strings.HasPrefix(token, query) {
			return true
		}

		var distance int
		hasUppercase := regexp.MustCompile("[A-Z]")
		if hasUppercase.MatchString(query) {
			distance = fuzzy.RankMatchNormalized(query, token)
		} else {
			distance = fuzzy.RankMatchNormalizedFold(query, token)
		}

		if distance == -1 {
			continue
		}

		if distance == 0 {
			return true
		}

		distancePercent := (float64(distance) / float64(len(token))) * 100
		if distancePercent < 50 {
			return true
		}
	}
	return false
}
