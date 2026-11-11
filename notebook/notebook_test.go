package notebook

import (
	"reflect"
	"sort"
	"testing"

	"github.com/msp301/graph"
)

func TestMatchedTags(t *testing.T) {
	g := graph.New()
	g.Add(1, "note", nil)
	g.Add(2, "note", nil)
	g.Add(3, "note", nil)

	g.Add(4, "tag", "#test")
	g.Add(5, "tag", "#LLM")
	g.Add(6, "tag", "#LargeLanguageModel")
	g.Add(7, "tag", "#LargeLanguageModels")

	book := &Notebook{Notes: g}

	got := book.MatchedTags("llm")
	want := []matchedTag{
		{Term: "llm", Distance: 1, Tag: "#LLM", Vertex: graph.Vertex{Id: 5, Label: "tag", Value: "#LLM"}},
		{Term: "llm", Distance: 16, Tag: "#LargeLanguageModel", Vertex: graph.Vertex{Id: 6, Label: "tag", Value: "#LargeLanguageModel"}},
		{Term: "llm", Distance: 17, Tag: "#LargeLanguageModels", Vertex: graph.Vertex{Id: 7, Label: "tag", Value: "#LargeLanguageModels"}},
	}

	if !reflect.DeepEqual(got, want) {
		t.Fatalf("Expected: %+v\nGot: %+v", want, got)
	}
}

func tagNote(g *graph.Graph, noteId uint64, tagIds ...uint64) {
	for i, tagId := range tagIds {
		g.Edge(noteId, tagId, "tag")
		for _, coOccurringTagId := range tagIds[i+1:] {
			g.Edge(tagId, coOccurringTagId, "tag")
		}
	}
}

func resultNoteIds(results []Result) []uint64 {
	ids := []uint64{}
	for _, result := range results {
		vertex := result.Value.(graph.Vertex)
		if vertex.Label == "note" {
			ids = append(ids, vertex.Id)
		}
	}
	sort.Slice(ids, func(i, j int) bool { return ids[i] < ids[j] })
	return ids
}

func resultTags(results []Result) []string {
	tags := []string{}
	for _, result := range results {
		vertex := result.Value.(graph.Vertex)
		if vertex.Label == "tag" {
			tags = append(tags, vertex.Value.(string))
		}
	}
	sort.Strings(tags)
	return tags
}

func assertSearchByTags(t *testing.T, book *Notebook, terms []string, wantNotes []uint64, wantTags []string) {
	t.Helper()
	got := book.SearchByTags(terms...)
	if gotNotes := resultNoteIds(got); !reflect.DeepEqual(gotNotes, wantNotes) {
		t.Errorf("SearchByTags(%q) notes\nExpected: %v\nGot: %v", terms, wantNotes, gotNotes)
	}
	if gotTags := resultTags(got); !reflect.DeepEqual(gotTags, wantTags) {
		t.Errorf("SearchByTags(%q) tags\nExpected: %v\nGot: %v", terms, wantTags, gotTags)
	}
}

const (
	noteTakingTag uint64 = iota + 1
	largeLanguageModelsTag
	todoTag
	llmTag
	literatureTag
)

const (
	noteTakingAndLargeLanguageModelsNote uint64 = iota + 10
	noteTakingAndTodoNote
	llmAndTodoNote
	noteTakingAndLiteratureNote
	llmAndLiteratureNote
)

func bookWhereTagsCoOccurAcrossDifferentNotes() *Notebook {
	g := graph.New()
	g.Add(noteTakingTag, "tag", "#NoteTaking")
	g.Add(largeLanguageModelsTag, "tag", "#LargeLanguageModels")
	g.Add(todoTag, "tag", "#TODO")
	g.Add(llmTag, "tag", "#LLM")
	g.Add(literatureTag, "tag", "#Literature")

	g.Add(noteTakingAndLargeLanguageModelsNote, "note", nil)
	g.Add(noteTakingAndTodoNote, "note", nil)
	g.Add(llmAndTodoNote, "note", nil)
	g.Add(noteTakingAndLiteratureNote, "note", nil)
	g.Add(llmAndLiteratureNote, "note", nil)

	tagNote(g, noteTakingAndLargeLanguageModelsNote, noteTakingTag, largeLanguageModelsTag)
	tagNote(g, noteTakingAndTodoNote, noteTakingTag, todoTag)
	tagNote(g, llmAndTodoNote, llmTag, todoTag)
	tagNote(g, noteTakingAndLiteratureNote, noteTakingTag, literatureTag)
	tagNote(g, llmAndLiteratureNote, llmTag, literatureTag)

	return &Notebook{Notes: g}
}

func TestMatchedTags_RecordsTagAgainstEveryTermItMatches(t *testing.T) {
	g := graph.New()
	g.Add(1, "tag", "#LargeLanguageModels")
	g.Add(2, "tag", "#LLM")

	book := &Notebook{Notes: g}

	type termTag struct{ Term, Tag string }
	got := map[termTag]bool{}
	for _, matched := range book.MatchedTags("llm", "large") {
		got[termTag{matched.Term, matched.Tag}] = true
	}
	want := map[termTag]bool{
		{"llm", "#LLM"}:                   true,
		{"llm", "#LargeLanguageModels"}:   true,
		{"large", "#LargeLanguageModels"}: true,
	}

	if !reflect.DeepEqual(got, want) {
		t.Fatalf("Expected: %+v\nGot: %+v", want, got)
	}
}

func TestSearchByTags_OneTagCanSatisfyMultipleTerms(t *testing.T) {
	const largeLanguageModels, llm uint64 = 1, 2
	const largeLanguageModelsOnlyNote, llmOnlyNote uint64 = 10, 11

	g := graph.New()
	g.Add(largeLanguageModels, "tag", "#LargeLanguageModels")
	g.Add(llm, "tag", "#LLM")
	g.Add(largeLanguageModelsOnlyNote, "note", nil)
	g.Add(llmOnlyNote, "note", nil)
	tagNote(g, largeLanguageModelsOnlyNote, largeLanguageModels)
	tagNote(g, llmOnlyNote, llm)

	book := &Notebook{Notes: g}

	for _, terms := range [][]string{{"llm", "large"}, {"large", "llm"}} {
		assertSearchByTags(t, book, terms, []uint64{largeLanguageModelsOnlyNote}, []string{"#LargeLanguageModels"})
	}
}

func TestSearchByTags_EveryTermMustMatchRegardlessOfOrder(t *testing.T) {
	const css, accessibilityWhichAlsoMatchesCss uint64 = 1, 2
	const cssOnlyNote, cssAndAccessibilityNote, accessibilityOnlyNote uint64 = 10, 11, 12

	g := graph.New()
	g.Add(css, "tag", "#CSS")
	g.Add(accessibilityWhichAlsoMatchesCss, "tag", "#Accessibility")
	g.Add(cssOnlyNote, "note", nil)
	g.Add(cssAndAccessibilityNote, "note", nil)
	g.Add(accessibilityOnlyNote, "note", nil)
	tagNote(g, cssOnlyNote, css)
	tagNote(g, cssAndAccessibilityNote, css, accessibilityWhichAlsoMatchesCss)
	tagNote(g, accessibilityOnlyNote, accessibilityWhichAlsoMatchesCss)

	book := &Notebook{Notes: g}

	for _, terms := range [][]string{{"accessibility", "css"}, {"css", "accessibility"}} {
		assertSearchByTags(t, book, terms,
			[]uint64{cssAndAccessibilityNote, accessibilityOnlyNote},
			[]string{"#Accessibility", "#CSS"})
	}
}

func TestSearchByTags_TermWithoutMatchingTagReturnsNothing(t *testing.T) {
	g := graph.New()
	g.Add(1, "tag", "#foo")
	g.Add(10, "note", nil)
	tagNote(g, 10, 1)

	book := &Notebook{Notes: g}

	for _, terms := range [][]string{{"foo", "zzzqqq"}, {"zzzqqq", "foo"}} {
		if got := book.SearchByTags(terms...); len(got) != 0 {
			t.Errorf("SearchByTags(%q)\nExpected no results\nGot: %+v", terms, got)
		}
	}
}

func TestSearchByTags_ListsMatchedTagsButNotTagsCoOccurringOnOtherNotes(t *testing.T) {
	book := bookWhereTagsCoOccurAcrossDifferentNotes()

	assertSearchByTags(t, book, []string{"llm", "notetaking"},
		[]uint64{noteTakingAndLargeLanguageModelsNote},
		[]string{"#LargeLanguageModels", "#NoteTaking"})
}

func TestSearchByTags_ListsEveryTagOfMatchedNotes(t *testing.T) {
	book := bookWhereTagsCoOccurAcrossDifferentNotes()

	assertSearchByTags(t, book, []string{"notetaking"},
		[]uint64{noteTakingAndLargeLanguageModelsNote, noteTakingAndTodoNote, noteTakingAndLiteratureNote},
		[]string{"#LargeLanguageModels", "#Literature", "#NoteTaking", "#TODO"})
}

func TestSearchByTags_ListsTagsBeforeNotes(t *testing.T) {
	book := bookWhereTagsCoOccurAcrossDifferentNotes()

	seenNote := false
	for _, result := range book.SearchByTags("notetaking") {
		vertex := result.Value.(graph.Vertex)
		if vertex.Label == "note" {
			seenNote = true
		} else if seenNote {
			t.Fatalf("Expected all tags before notes, got %q after a note", vertex.Value)
		}
	}
}

func TestSearchByTags_AddingAListedTagRefinesResults(t *testing.T) {
	book := bookWhereTagsCoOccurAcrossDifferentNotes()

	for _, terms := range [][]string{{"notetaking"}, {"llm"}, {"llm", "notetaking"}, {"notetaking", "todo"}} {
		results := book.SearchByTags(terms...)
		notes := map[uint64]bool{}
		for _, id := range resultNoteIds(results) {
			notes[id] = true
		}
		if len(notes) == 0 {
			t.Fatalf("SearchByTags(%q) expected results", terms)
		}

		for _, tag := range resultTags(results) {
			refinedTerms := append(append([]string{}, terms...), tag)
			refined := resultNoteIds(book.SearchByTags(refinedTerms...))
			if len(refined) == 0 {
				t.Errorf("SearchByTags(%q) listed %s, but adding it returned no notes", terms, tag)
			}
			for _, id := range refined {
				if !notes[id] {
					t.Errorf("SearchByTags(%q) returned note %d, which SearchByTags(%q) did not", refinedTerms, id, terms)
				}
			}
		}
	}
}

func TestSearchByTag(t *testing.T) {
	g := graph.New()
	g.Add(1, "note", nil)
	g.Add(2, "note", nil)
	g.Add(3, "tag", "#foo")
	g.Add(4, "note", nil)
	g.Add(5, "tag", "#bar")
	g.Edge(2, 3, "tag")
	g.Edge(1, 5, "tag")
	book := &Notebook{Notes: g}

	got := book.SearchByTags("#foo")
	want := []Result{
		{Line: -1, Value: graph.Vertex{Id: 3, Label: "tag", Value: "#foo"}},
		{Line: -1, Value: graph.Vertex{Id: 2, Label: "note"}},
	}

	if !reflect.DeepEqual(got, want) {
		t.Fatalf("Expected: %+v\nGot: %+v", want, got)
	}
}

func TestSearchByTags_PrioritisesMatches(t *testing.T) {
	g := graph.New()
	g.Add(1, "note", nil)
	g.Add(2, "note", nil)
	g.Add(3, "note", nil)

	g.Add(4, "tag", "#test")
	g.Add(5, "tag", "#LLM")
	g.Add(6, "tag", "#LargeLanguageModel")
	g.Add(7, "tag", "#LargeLanguageModels")

	g.Edge(2, 4, "tag")
	g.Edge(3, 4, "tag")
	g.Edge(3, 7, "tag")

	book := &Notebook{Notes: g}

	got := book.SearchByTags("test", "llm")
	want := []Result{
		{Line: -1, Value: graph.Vertex{Id: 4, Label: "tag", Value: "#test"}},
		{Line: -1, Value: graph.Vertex{Id: 7, Label: "tag", Value: "#LargeLanguageModels"}},
		{Line: -1, Value: graph.Vertex{Id: 3, Label: "note"}},
	}

	if !reflect.DeepEqual(got, want) {
		t.Fatalf("Expected: %+v\nGot: %+v", want, got)
	}
}

func TestTagIntersection(t *testing.T) {
	g := graph.New()
	g.Add(1, "note", nil)
	g.Add(2, "note", nil)
	g.Add(3, "note", nil)

	g.Add(4, "tag", "#test")
	g.Add(5, "tag", "#LLM")
	g.Add(6, "tag", "#LargeLanguageModel")
	g.Add(7, "tag", "#LargeLanguageModels")

	g.Edge(2, 4, "tag")
	g.Edge(3, 4, "tag")
	g.Edge(3, 7, "tag")

	book := &Notebook{Notes: g}
	tags := []matchedTag{
		{Term: "test", Distance: 1, Tag: "#test", Vertex: graph.Vertex{Id: 4, Label: "tag", Value: "#test"}},
		{Term: "llm", Distance: 1, Tag: "#LLM", Vertex: graph.Vertex{Id: 5, Label: "tag", Value: "#LLM"}},
		{Term: "llm", Distance: 16, Tag: "#LargeLanguageModel", Vertex: graph.Vertex{Id: 6, Label: "tag", Value: "#LargeLanguageModel"}},
		{Term: "llm", Distance: 17, Tag: "#LargeLanguageModels", Vertex: graph.Vertex{Id: 7, Label: "tag", Value: "#LargeLanguageModels"}},
	}

	got := book.TagIntersection(tags)
	want := []graph.Vertex{{Id: 3, Label: "note"}}

	if !reflect.DeepEqual(got, want) {
		t.Fatalf("Expected: %+v\nGot: %+v", want, got)
	}
}

func TestTagIntersection_MatchesNoteWithSecondClosestTagForTerm(t *testing.T) {
	g := graph.New()
	g.Add(1, "note", nil)
	g.Add(2, "note", nil)

	g.Add(3, "tag", "#foo")
	g.Add(4, "tag", "#bar")
	g.Add(5, "tag", "#foobar")

	g.Edge(1, 3, "tag")
	g.Edge(1, 5, "tag")
	g.Edge(2, 4, "tag")

	book := &Notebook{Notes: g}
	tags := []matchedTag{
		{Term: "foo", Distance: 1, Tag: "#foo", Vertex: graph.Vertex{Id: 3, Label: "tag", Value: "#foo"}},
		{Term: "bar", Distance: 1, Tag: "#bar", Vertex: graph.Vertex{Id: 4, Label: "tag", Value: "#bar"}},
		{Term: "bar", Distance: 4, Tag: "#foobar", Vertex: graph.Vertex{Id: 5, Label: "tag", Value: "#foobar"}},
	}

	got := book.TagIntersection(tags)
	want := []graph.Vertex{{Id: 1, Label: "note"}}

	if !reflect.DeepEqual(got, want) {
		t.Fatalf("Expected: %+v\nGot: %+v", want, got)
	}
}

func TestTagIntersection_DoesNotCountSameSearchTermTwice(t *testing.T) {
	g := graph.New()
	g.Add(1, "note", nil) // has #golang + #gopher (both match "go", no "web")
	g.Add(2, "note", nil) // has #golang + #web (correct match)

	g.Add(3, "tag", "#golang")
	g.Add(4, "tag", "#gopher")
	g.Add(5, "tag", "#web")

	g.Edge(1, 3, "tag")
	g.Edge(1, 4, "tag")
	g.Edge(2, 3, "tag")
	g.Edge(2, 5, "tag")

	book := &Notebook{Notes: g}
	tags := []matchedTag{
		{Term: "go", Distance: 1, Tag: "#golang", Vertex: graph.Vertex{Id: 3, Label: "tag", Value: "#golang"}},
		{Term: "go", Distance: 2, Tag: "#gopher", Vertex: graph.Vertex{Id: 4, Label: "tag", Value: "#gopher"}},
		{Term: "web", Distance: 1, Tag: "#web", Vertex: graph.Vertex{Id: 5, Label: "tag", Value: "#web"}},
	}

	got := book.TagIntersection(tags)
	want := []graph.Vertex{{Id: 2, Label: "note"}}

	if !reflect.DeepEqual(got, want) {
		t.Fatalf("Expected: %+v\nGot: %+v", want, got)
	}
}
