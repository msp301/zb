package browser

import (
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/msp301/zb/bookshelf"
	"github.com/msp301/zb/notebook"
)

var fixtures = map[string]string{
	"202401010001.md": "202401010001\n# Graph theory\n\nVertices and edges #maths\n\nSee [[202401010002]]\n",
	"202401010002.md": "202401010002\n# Zettelkasten method\n\nLinking notes builds a graph of ideas #notes\n",
	"202401010003.md": "202401010003\n# Sourdough starter\n\nFeed the starter daily #baking\n",
	"202301010001.md": "202301010001\n# Schema diagrams\n\nEntity relationships #diagrams #database\n",
	"202301010002.md": "202301010002\n# Sequence diagrams\n\nMessages over time #diagrams\n",
	"202301010003.md": "202301010003\n# Index tuning\n\nCovering indexes #database\n",
}

func newTestModel(t *testing.T, query string) *model {
	t.Helper()

	dir := t.TempDir()
	for name, content := range fixtures {
		if err := os.WriteFile(filepath.Join(dir, name), []byte(content), 0666); err != nil {
			t.Fatal(err)
		}
	}

	options := Options{
		Dirs:  []string{dir},
		Query: query,
		Load: func() (*notebook.Notebook, error) {
			err, book := bookshelf.Read([]string{dir})
			return book, err
		},
	}

	book, err := options.Load()
	if err != nil {
		t.Fatal(err)
	}

	m := newModel(options, book)
	m.delay = 0
	send(m, tea.WindowSizeMsg{Width: 120, Height: 30})

	return m
}

func send(m *model, msg tea.Msg) {
	_, cmd := m.Update(msg)
	run(m, cmd)
}

func run(m *model, cmd tea.Cmd) {
	if cmd == nil {
		return
	}

	switch msg := cmd().(type) {
	case tea.BatchMsg:
		for _, batched := range msg {
			run(m, batched)
		}
	case searchDueMsg, searchDoneMsg, reloadedMsg:
		send(m, msg)
	}
}

func typeText(m *model, text string) {
	for _, r := range text {
		send(m, tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{r}})
	}
}

func press(m *model, key tea.KeyType) {
	send(m, tea.KeyMsg{Type: key})
}

func titles(m *model) []string {
	var titles []string
	for _, r := range m.rows {
		titles = append(titles, r.title)
	}
	return titles
}

func expectFacets(t *testing.T, m *model, want ...string) {
	t.Helper()
	var got []string
	for _, f := range m.facets {
		got = append(got, fmt.Sprintf("%s %d", f.tag, f.count))
	}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("Expected tags: %v\nGot: %v", want, got)
	}
}

func expectTitles(t *testing.T, m *model, want ...string) {
	t.Helper()
	if got := titles(m); !reflect.DeepEqual(got, want) {
		t.Fatalf("Expected rows: %v\nGot: %v", want, got)
	}
}

func TestEmptyQueryListsRecentNotes(t *testing.T) {
	m := newTestModel(t, "")

	expectTitles(t, m, "Sourdough starter", "Zettelkasten method", "Graph theory", "Index tuning", "Sequence diagrams", "Schema diagrams")
}

func TestSelectionWrapsAroundTheList(t *testing.T) {
	m := newTestModel(t, "#diagrams")

	press(m, tea.KeyUp)

	if selected, _ := m.selected(); selected.title != "Sequence diagrams" {
		t.Fatalf("Expected up on the first row to select the last, got %q", selected.title)
	}

	press(m, tea.KeyDown)

	if selected, _ := m.selected(); selected.title != "Schema diagrams" {
		t.Fatalf("Expected down on the last row to select the first, got %q", selected.title)
	}
}

func TestTypingSearchesNoteText(t *testing.T) {
	m := newTestModel(t, "")

	typeText(m, "starter")

	expectTitles(t, m, "Sourdough starter")
	if !strings.Contains(m.View(), "Feed the starter daily") {
		t.Fatalf("Expected preview to show the selected note, got:\n%s", m.View())
	}
}

func TestInitialQueryIsSearched(t *testing.T) {
	m := newTestModel(t, "starter")

	expectTitles(t, m, "Sourdough starter")
}

func TestHashPrefixSearchesByTag(t *testing.T) {
	m := newTestModel(t, "")

	typeText(m, "#baking")

	expectTitles(t, m, "Sourdough starter")
}

func TestTabShowsRelatedAndShiftTabReturns(t *testing.T) {
	m := newTestModel(t, "")
	press(m, tea.KeyDown)
	press(m, tea.KeyDown)

	press(m, tea.KeyTab)

	if got := m.current().label(); got != "related: Graph theory" {
		t.Fatalf("Expected to be viewing related notes, got %q", got)
	}
	expectTitles(t, m, "Zettelkasten method")
	expectFacets(t, m, "#maths 1")

	press(m, tea.KeyShiftTab)

	expectTitles(t, m, "Sourdough starter", "Zettelkasten method", "Graph theory", "Index tuning", "Sequence diagrams", "Schema diagrams")
	if selected, _ := m.selected(); selected.title != "Graph theory" {
		t.Fatalf("Expected selection to be restored, got %q", selected.title)
	}
}

func TestTagSearchListsNotesAndShowsOtherTagsAsFacets(t *testing.T) {
	m := newTestModel(t, "#diagrams")

	expectTitles(t, m, "Schema diagrams", "Sequence diagrams")
	expectFacets(t, m, "#database 1")
	if !strings.Contains(m.View(), "#database 1") {
		t.Fatalf("Expected related tags under the query, got:\n%s", m.View())
	}
}

func TestTypingHashSuggestsTheMostUsedTag(t *testing.T) {
	m := newTestModel(t, "")

	typeText(m, "#")

	if m.suggestion != "#database" {
		t.Fatalf("Expected the most used tag to be suggested, got %q", m.suggestion)
	}
	expectFacets(t, m, "#database 2", "#diagrams 2", "#baking 1", "#maths 1", "#notes 1")
	expectTitles(t, m, "Schema diagrams", "Index tuning")
}

func TestTabCompletesTheSuggestedTag(t *testing.T) {
	m := newTestModel(t, "")
	typeText(m, "#di")

	press(m, tea.KeyTab)

	if got := m.input.Value(); got != "#diagrams " {
		t.Fatalf("Expected the tag to be completed, got %q", got)
	}
	if len(m.history) != 1 {
		t.Fatalf("Expected completing the first tag to stay in place, got %d entries", len(m.history))
	}
	expectTitles(t, m, "Schema diagrams", "Sequence diagrams")
}

func TestTabNarrowsTagSearchWithSuggestedTag(t *testing.T) {
	m := newTestModel(t, "#diagrams")
	typeText(m, " #da")

	if m.suggestion != "#database" {
		t.Fatalf("Expected a tag of the listed notes to be suggested, got %q", m.suggestion)
	}
	expectTitles(t, m, "Schema diagrams")

	press(m, tea.KeyTab)

	if got := m.input.Value(); got != "#diagrams #database " {
		t.Fatalf("Expected the tag to be added to the query, got %q", got)
	}
	expectTitles(t, m, "Schema diagrams")

	press(m, tea.KeyShiftTab)

	if got := m.input.Value(); got != "#diagrams" {
		t.Fatalf("Expected to return to the broader search, got %q", got)
	}
	expectTitles(t, m, "Schema diagrams", "Sequence diagrams")
}

func TestEnterSearchesOnlyTheSuggestedTag(t *testing.T) {
	m := newTestModel(t, "#diagrams")
	typeText(m, " #da")

	press(m, tea.KeyEnter)

	if got := m.input.Value(); got != "#database " {
		t.Fatalf("Expected query to be the tag alone, got %q", got)
	}
	expectTitles(t, m, "Schema diagrams", "Index tuning")

	press(m, tea.KeyShiftTab)

	if got := m.input.Value(); got != "#diagrams" {
		t.Fatalf("Expected to return to the previous search, got %q", got)
	}
}

func TestSuggestionsOnlyOfferTagsOfTheListedNotes(t *testing.T) {
	m := newTestModel(t, "#diagrams")

	typeText(m, " #")

	expectFacets(t, m, "#database 1")
}

func TestCompleteTagIsNotSuggested(t *testing.T) {
	m := newTestModel(t, "#diagrams")

	press(m, tea.KeyTab)

	if got := m.current().label(); got != "related: Schema diagrams" {
		t.Fatalf("Expected Tab to show related notes, got %q", got)
	}
}

func TestTagsOfTheRelatedNoteAreSuggestedFirst(t *testing.T) {
	m := newTestModel(t, "starter")
	press(m, tea.KeyTab)

	typeText(m, "#")

	if m.suggestion != "#baking" {
		t.Fatalf("Expected the note's own tag to be suggested, got %q", m.suggestion)
	}
}

func TestTextAndTagsSearchTogether(t *testing.T) {
	m := newTestModel(t, "time #diagrams")

	expectTitles(t, m, "Sequence diagrams")

	m = newTestModel(t, "starter #diagrams")

	expectTitles(t, m)
}

func TestTypingInRelatedViewStartsNewSearch(t *testing.T) {
	m := newTestModel(t, "")
	press(m, tea.KeyTab)

	typeText(m, "graph")
	press(m, tea.KeyShiftTab)

	if got := m.current().label(); got != "related: Sourdough starter" {
		t.Fatalf("Expected to return to the related view, got %q", got)
	}
}

func TestEscapeClearsThenGoesBackThenQuits(t *testing.T) {
	m := newTestModel(t, "starter")
	press(m, tea.KeyTab)

	press(m, tea.KeyEsc)
	if len(m.history) != 1 || m.input.Value() != "starter" {
		t.Fatalf("Expected escape to go back to the search, got %d entries and query %q", len(m.history), m.input.Value())
	}

	press(m, tea.KeyEsc)
	if m.input.Value() != "" {
		t.Fatalf("Expected escape to clear the query, got %q", m.input.Value())
	}

	_, cmd := m.Update(tea.KeyMsg{Type: tea.KeyEsc})
	if _, ok := cmd().(tea.QuitMsg); !ok {
		t.Fatal("Expected escape on an empty root query to quit")
	}
}

func TestNotebookChangesAppearAfterEditing(t *testing.T) {
	m := newTestModel(t, "")
	path := filepath.Join(m.options.Dirs[0], "202401010004.md")
	if err := os.WriteFile(path, []byte("202401010004\n# Brand new note\n"), 0666); err != nil {
		t.Fatal(err)
	}

	send(m, editorClosedMsg{})

	if got := titles(m); !contains(got, "Brand new note") {
		t.Fatalf("Expected new note after reload, got %v", got)
	}
}

func TestNarrowTerminalStacksListAbovePreview(t *testing.T) {
	m := newTestModel(t, "starter")

	send(m, tea.WindowSizeMsg{Width: 50, Height: 20})

	lines := strings.Split(m.View(), "\n")
	if len(lines) != 20 {
		t.Fatalf("Expected view to fill 20 lines, got %d", len(lines))
	}
	if !strings.Contains(m.View(), "Feed the starter daily") {
		t.Fatalf("Expected preview below the list, got:\n%s", m.View())
	}
}

func contains(values []string, want string) bool {
	for _, value := range values {
		if value == want {
			return true
		}
	}
	return false
}

func TestRelatedViewExcludesTheNoteItself(t *testing.T) {
	m := newTestModel(t, "")

	press(m, tea.KeyTab)

	if got := titles(m); contains(got, "Sourdough starter") {
		t.Fatalf("Expected the note itself to be left out, got %v", got)
	}
}

func wheel(m *model, button tea.MouseButton, x int, y int) {
	send(m, tea.MouseMsg{Action: tea.MouseActionPress, Button: button, X: x, Y: y})
}

func TestWheelOverPreviewScrollsTheNoteAndKeepsTheSelection(t *testing.T) {
	m := newTestModel(t, "")
	send(m, tea.WindowSizeMsg{Width: 120, Height: 8})
	selected, _ := m.selected()

	wheel(m, tea.MouseButtonWheelDown, 100, 4)

	if m.preview.YOffset == 0 {
		t.Fatal("expected the preview to scroll down")
	}
	if m.listOffset != 0 {
		t.Fatalf("expected the list to stay put, got offset %d", m.listOffset)
	}
	if after, _ := m.selected(); after.note != selected.note {
		t.Fatal("expected the selection to stay on the same note")
	}

	wheel(m, tea.MouseButtonWheelUp, 100, 4)

	if m.preview.YOffset != 0 {
		t.Fatalf("expected the preview to scroll back to the top, got offset %d", m.preview.YOffset)
	}
}

func TestWheelOverListScrollsTheListAndKeepsTheSelectionVisible(t *testing.T) {
	m := newTestModel(t, "")
	send(m, tea.WindowSizeMsg{Width: 120, Height: 8})

	wheel(m, tea.MouseButtonWheelDown, 5, 4)

	if m.listOffset != 2 {
		t.Fatalf("expected the list to scroll to its end, got offset %d", m.listOffset)
	}
	if m.current().cursor != 2 {
		t.Fatalf("expected the selection to follow the list, got %d", m.current().cursor)
	}
	if m.preview.YOffset != 0 {
		t.Fatalf("expected the preview to stay at the top, got offset %d", m.preview.YOffset)
	}

	wheel(m, tea.MouseButtonWheelUp, 5, 4)

	if m.listOffset != 0 {
		t.Fatalf("expected the list to scroll back to the top, got offset %d", m.listOffset)
	}
}

func TestWheelInStackedLayoutUsesThePaneUnderTheCursor(t *testing.T) {
	m := newTestModel(t, "")
	send(m, tea.WindowSizeMsg{Width: 60, Height: 13})

	wheel(m, tea.MouseButtonWheelDown, 5, 3)

	if m.listOffset != 3 || m.preview.YOffset != 0 {
		t.Fatalf("expected only the list to scroll, got list %d preview %d", m.listOffset, m.preview.YOffset)
	}

	wheel(m, tea.MouseButtonWheelDown, 5, 8)

	if m.listOffset != 3 || m.preview.YOffset == 0 {
		t.Fatalf("expected only the preview to scroll, got list %d preview %d", m.listOffset, m.preview.YOffset)
	}
}

func click(m *model, x int, y int) {
	send(m, tea.MouseMsg{Action: tea.MouseActionPress, Button: tea.MouseButtonLeft, X: x, Y: y})
}

func TestClickingARowSelectsIt(t *testing.T) {
	m := newTestModel(t, "")
	want := m.rows[2].note

	click(m, 5, 5)

	if m.current().cursor != 2 {
		t.Fatalf("expected the third row to be selected, got %d", m.current().cursor)
	}
	if !strings.Contains(m.preview.View(), want.Title) {
		t.Fatalf("expected the preview to show %q", want.Title)
	}
}

func TestClickingAScrolledListSelectsTheRowUnderTheCursor(t *testing.T) {
	m := newTestModel(t, "")
	send(m, tea.WindowSizeMsg{Width: 120, Height: 8})
	wheel(m, tea.MouseButtonWheelDown, 5, 4)

	click(m, 5, 6)

	if m.current().cursor != 5 {
		t.Fatalf("expected the last row to be selected, got %d", m.current().cursor)
	}
}

func TestClickingOutsideTheRowsKeepsTheSelection(t *testing.T) {
	m := newTestModel(t, "")

	click(m, 5, 20)
	click(m, 100, 5)
	click(m, 5, 0)

	if m.current().cursor != 0 {
		t.Fatalf("expected the selection to stay on the first row, got %d", m.current().cursor)
	}
}
