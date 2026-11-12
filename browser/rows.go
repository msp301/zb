package browser

import (
	"sort"
	"strings"

	"github.com/msp301/graph"
	"github.com/msp301/zb"
	"github.com/msp301/zb/notebook"
	"github.com/msp301/zb/util"
)

const fuzzyMinLength = 3

type row struct {
	title   string
	context string
	line    int
	note    *zb.Note
}

type facet struct {
	tag   string
	count int
}

func lookup(book *notebook.Notebook, loc location) []row {
	rows := toRows(results(book, loc))
	if loc.note == nil {
		return rows
	}

	related := rows[:0]
	for _, r := range rows {
		if r.note.Id != loc.note.Id {
			related = append(related, r)
		}
	}
	return related
}

func results(book *notebook.Notebook, loc location) []notebook.Result {
	if loc.note != nil {
		return book.SearchRelated(loc.note.Id)
	}

	tags, text := splitTerms(strings.Fields(loc.query))
	switch {
	case len(tags) == 0 && len(text) == 0:
		return book.Recent(recentLimit)
	case len(text) == 0:
		return book.SearchByTags(tags...)
	case len(tags) == 0:
		return book.Search(text...)
	}

	tagged := make(map[uint64]bool)
	for _, result := range book.SearchByTags(tags...) {
		if vertex, ok := result.Value.(graph.Vertex); ok {
			tagged[vertex.Id] = true
		}
	}

	var matches []notebook.Result
	for _, result := range book.Search(text...) {
		if vertex, ok := result.Value.(graph.Vertex); ok && tagged[vertex.Id] {
			matches = append(matches, result)
		}
	}
	return matches
}

func splitTerms(terms []string) (tags []string, text []string) {
	for _, term := range terms {
		if strings.HasPrefix(term, "#") {
			tags = append(tags, term)
		} else {
			text = append(text, term)
		}
	}
	return tags, text
}

func splitPartial(query string) (base []string, partial string) {
	terms := strings.Fields(query)
	if len(terms) == 0 || strings.TrimRight(query, " \t") != query {
		return terms, ""
	}

	last := terms[len(terms)-1]
	if !strings.HasPrefix(last, "#") {
		return terms, ""
	}
	return terms[:len(terms)-1], last
}

func toRows(results []notebook.Result) []row {
	rows := make([]row, 0, len(results))
	for _, result := range results {
		vertex, ok := result.Value.(graph.Vertex)
		if !ok {
			continue
		}

		note, ok := vertex.Value.(zb.Note)
		if !ok {
			continue
		}

		context := strings.Join(strings.Fields(result.Context), " ")
		line := result.Line
		if context == "" {
			line = note.Start
		}
		if strings.TrimLeft(context, "# ") == note.Title {
			context = ""
		}
		rows = append(rows, row{title: note.Title, context: context, line: line, note: &note})
	}
	return rows
}

func allNotes(book *notebook.Notebook) []zb.Note {
	var notes []zb.Note
	for _, vertex := range book.Notes.Vertices {
		if note, ok := vertex.Value.(zb.Note); ok {
			notes = append(notes, note)
		}
	}
	return notes
}

func rowNotes(rows []row) []zb.Note {
	seen := make(map[uint64]bool)
	var notes []zb.Note
	for _, r := range rows {
		if !seen[r.note.Id] {
			seen[r.note.Id] = true
			notes = append(notes, *r.note)
		}
	}
	return notes
}

func tagCounts(notes []zb.Note) map[string]int {
	counts := make(map[string]int)
	for _, note := range notes {
		for _, tag := range note.Tags {
			counts[tag]++
		}
	}
	return counts
}

func rank(counts map[string]int, exclude []string) []facet {
	facets := make([]facet, 0, len(counts))
	for tag, count := range counts {
		if !containsFold(exclude, tag) {
			facets = append(facets, facet{tag: tag, count: count})
		}
	}
	sort.Slice(facets, func(i, j int) bool {
		if facets[i].count != facets[j].count {
			return facets[i].count > facets[j].count
		}
		return strings.ToLower(facets[i].tag) < strings.ToLower(facets[j].tag)
	})
	return facets
}

func complete(candidates []facet, partial string, preferred []string) []facet {
	var matches []facet
	for _, candidate := range candidates {
		fuzzy := len(partial) >= fuzzyMinLength && util.Matches(partial, candidate.tag)
		if hasPrefixFold(candidate.tag, partial) || fuzzy {
			matches = append(matches, candidate)
		}
	}

	score := func(f facet) int {
		score := 0
		if hasPrefixFold(f.tag, partial) {
			score += 2
		}
		if containsFold(preferred, f.tag) {
			score++
		}
		return score
	}
	sort.SliceStable(matches, func(i, j int) bool {
		return score(matches[i]) > score(matches[j])
	})
	return matches
}

func hasPrefixFold(text string, prefix string) bool {
	return strings.HasPrefix(strings.ToLower(text), strings.ToLower(prefix))
}

func containsFold(values []string, want string) bool {
	for _, value := range values {
		if strings.EqualFold(value, want) {
			return true
		}
	}
	return false
}
