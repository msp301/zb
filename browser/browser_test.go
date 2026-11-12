package browser

import (
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
	if got := titles(m); !contains(got, "Zettelkasten method") || !contains(got, "#maths") {
		t.Fatalf("Expected linked note and tag, got %v", got)
	}

	press(m, tea.KeyShiftTab)

	expectTitles(t, m, "Sourdough starter", "Zettelkasten method", "Graph theory", "Index tuning", "Sequence diagrams", "Schema diagrams")
	if selected, _ := m.selected(); selected.title != "Graph theory" {
		t.Fatalf("Expected selection to be restored, got %q", selected.title)
	}
}

func TestTagSearchHidesTheSearchedTag(t *testing.T) {
	m := newTestModel(t, "#diagrams")

	expectTitles(t, m, "#database", "Schema diagrams", "Sequence diagrams")
}

func TestTabOnTagNarrowsTagSearch(t *testing.T) {
	m := newTestModel(t, "#diagrams")

	press(m, tea.KeyTab)

	if got := m.input.Value(); got != "#diagrams #database" {
		t.Fatalf("Expected the tag to be added to the query, got %q", got)
	}
	expectTitles(t, m, "Schema diagrams")

	press(m, tea.KeyShiftTab)

	expectTitles(t, m, "#database", "Schema diagrams", "Sequence diagrams")
}

func TestEnterOnTagSearchesOnlyThatTag(t *testing.T) {
	m := newTestModel(t, "#diagrams")

	press(m, tea.KeyEnter)

	if got := m.input.Value(); got != "#database" {
		t.Fatalf("Expected query to be the tag, got %q", got)
	}
	if len(m.history) != 2 {
		t.Fatalf("Expected a new history entry, got %d", len(m.history))
	}
	expectTitles(t, m, "#diagrams", "Schema diagrams", "Index tuning")
}

func TestTabOnTagOfRelatedNoteSearchesThatTag(t *testing.T) {
	m := newTestModel(t, "starter")
	press(m, tea.KeyTab)

	press(m, tea.KeyTab)

	if got := m.input.Value(); got != "#baking" {
		t.Fatalf("Expected query to be the tag, got %q", got)
	}
}

func TestTagPreviewCountsNotesWithinTagSearch(t *testing.T) {
	m := newTestModel(t, "#diagrams")

	if !strings.Contains(m.View(), "1 of 2 notes here") {
		t.Fatalf("Expected preview to count the notes narrowing would keep, got:\n%s", m.View())
	}
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
