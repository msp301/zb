package browser

import (
	"sort"
	"strings"

	"github.com/msp301/graph"
	"github.com/msp301/zb"
	"github.com/msp301/zb/notebook"
)

type row struct {
	title   string
	context string
	line    int
	note    *zb.Note
	tag     string
	tagId   uint64
}

func lookup(book *notebook.Notebook, loc location) []row {
	rows := toRows(results(book, loc))
	if loc.note == nil {
		return rows
	}

	related := rows[:0]
	for _, r := range rows {
		if r.note == nil || r.note.Id != loc.note.Id {
			related = append(related, r)
		}
	}
	return related
}

func results(book *notebook.Notebook, loc location) []notebook.Result {
	if loc.note != nil {
		return book.SearchRelated(loc.note.Id)
	}

	terms := strings.Fields(loc.query)
	if len(terms) == 0 {
		return book.Recent(recentLimit)
	}

	if isTagQuery(terms) {
		return book.SearchByTags(terms...)
	}

	return book.Search(terms...)
}

func isTagQuery(terms []string) bool {
	for _, term := range terms {
		if !strings.HasPrefix(term, "#") {
			return false
		}
	}
	return true
}

func toRows(results []notebook.Result) []row {
	rows := make([]row, 0, len(results))
	for _, result := range results {
		vertex, ok := result.Value.(graph.Vertex)
		if !ok {
			continue
		}

		switch value := vertex.Value.(type) {
		case zb.Note:
			note := value
			context := strings.Join(strings.Fields(result.Context), " ")
			line := result.Line
			if context == "" {
				line = note.Start
			}
			if strings.TrimLeft(context, "# ") == note.Title {
				context = ""
			}
			rows = append(rows, row{title: note.Title, context: context, line: line, note: &note})
		case string:
			rows = append(rows, row{title: value, tag: value, tagId: vertex.Id})
		}
	}
	return rows
}

func taggedNotes(book *notebook.Notebook, tagId uint64) []zb.Note {
	var notes []zb.Note
	for id := range book.Notes.Adjacency[tagId] {
		if note, ok := book.Notes.Vertices[id].Value.(zb.Note); ok {
			notes = append(notes, note)
		}
	}
	sort.Slice(notes, func(i, j int) bool {
		return notes[i].Id > notes[j].Id
	})
	return notes
}
