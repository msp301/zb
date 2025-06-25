package notebook

import (
	"github.com/msp301/graph"
	"github.com/msp301/zb"
	"github.com/msp301/zb/util"
	"strings"
)

func (book *Notebook) Search(query ...string) []Result {
	results := []Result{}
	traversal := graph.Traversal(book.Notes)
	for vertex := range traversal.V().HasLabel("note").Iterate() {
		var context []util.ContextMatch
		matched := false
		startLine := 0
		switch val := vertex.Value.(type) {
		case zb.Note:
			content := val.Content
			paragraphs := extractParagraphs(content)
			startLine = val.Start
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
				Context: context.Text,
				Line:    startLine + context.Line - 1,
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
