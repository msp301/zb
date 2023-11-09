package notebook

import (
	"github.com/msp301/zb/graph"
	"github.com/msp301/zb/parser"
	"github.com/msp301/zb/util"
	"strings"
)

func (book *Notebook) Search(query ...string) []Result {
	results := []Result{}
	traversal := graph.Traversal(book.Notes)
	for vertex := range traversal.V().HasLabel("note").Iterate() {
		var context []string
		matched := false
		switch val := vertex.Properties["Value"].(type) {
		case parser.Note:
			content := val.Title + "\n\n" + val.Content
			paragraphs := extractParagraphs(content)
		PARAGRAPH:
			for _, paragraph := range paragraphs {
				termsMatched := 0
				for _, term := range query {
					if util.ParagraphMatches(paragraph, term) {
						termsMatched++
					}
				}

				if termsMatched != len(query) {
					continue PARAGRAPH
				}

				for _, term := range query {
					extracted, ok := util.ContextFold(paragraph, term)
					if ok {
						context = append(context, extracted...)
						matched = true
						break
					}
				}
				break
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
