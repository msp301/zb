package browser

import (
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/cursor"
	"github.com/charmbracelet/bubbles/textinput"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/msp301/zb"
	"github.com/msp301/zb/editor"
	"github.com/msp301/zb/notebook"
)

const recentLimit = 100
const searchDelay = 120 * time.Millisecond
const wheelStep = 3

type Options struct {
	Dirs  []string
	Query string
	Load  func() (*notebook.Notebook, error)
}

type location struct {
	query  string
	note   *zb.Note
	cursor int
}

func (loc location) label() string {
	if loc.note != nil {
		return "related: " + loc.note.Title
	}
	if query := strings.TrimSpace(loc.query); query != "" {
		return query
	}
	return "recent"
}

type searchDueMsg struct{ generation int }

type searchDoneMsg struct {
	generation int
	rows       []row
}

type editorClosedMsg struct{ err error }

type reloadedMsg struct {
	book *notebook.Notebook
	err  error
}

type model struct {
	options    Options
	book       *notebook.Notebook
	input      textinput.Model
	preview    viewport.Model
	history    []location
	rows       []row
	facets     []facet
	suggestion string
	tagTotals  map[string]int
	knownTags  map[string]bool
	baseQuery  string
	baseCounts map[string]int
	listOffset int
	width      int
	height     int
	status     string
	delay      time.Duration
	generation int
	applied    int
	searching  bool
	pending    bool
}

func Run(options Options) error {
	book, err := options.Load()
	if err != nil {
		return err
	}

	_, err = tea.NewProgram(newModel(options, book), tea.WithAltScreen(), tea.WithMouseCellMotion()).Run()
	return err
}

func newModel(options Options, book *notebook.Notebook) *model {
	input := textinput.New()
	input.Prompt = "> "
	input.Cursor.SetMode(cursor.CursorStatic)
	input.SetValue(options.Query)
	input.CursorEnd()
	input.ShowSuggestions = true
	input.Focus()

	m := &model{
		options: options,
		book:    book,
		input:   input,
		preview: viewport.New(0, 0),
		history: []location{{query: options.Query}},
		delay:   searchDelay,
	}
	m.indexTags()
	m.syncInput()
	m.syncFacets()
	m.rows = lookup(book, m.target())
	m.syncFacets()

	return m
}

func (m *model) Init() tea.Cmd {
	return nil
}

func (m *model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		m.resize()
		return m, nil

	case searchDueMsg:
		if msg.generation != m.generation {
			return m, nil
		}
		if m.searching {
			m.pending = true
			return m, nil
		}
		return m, m.search()

	case searchDoneMsg:
		m.searching = false
		if msg.generation == m.generation {
			m.rows = msg.rows
			m.applied = msg.generation
			m.syncFacets()
			m.clampCursor()
			m.showSelected()
		}
		if m.pending {
			m.pending = false
			return m, m.search()
		}
		return m, nil

	case editorClosedMsg:
		if msg.err != nil {
			m.status = "Editor error: " + msg.err.Error()
		}
		load := m.options.Load
		return m, func() tea.Msg {
			book, err := load()
			return reloadedMsg{book: book, err: err}
		}

	case reloadedMsg:
		if msg.err != nil {
			m.status = "Failed to reload notebook: " + msg.err.Error()
			return m, nil
		}
		m.book = msg.book
		m.indexTags()
		m.syncFacets()
		return m, m.refresh(0)

	case tea.KeyMsg:
		return m, m.handleKey(msg)

	case tea.MouseMsg:
		m.handleMouse(msg)
		return m, nil
	}

	return m, nil
}

func (m *model) handleMouse(msg tea.MouseMsg) {
	delta := 0
	switch msg.Button {
	case tea.MouseButtonWheelUp:
		delta = -wheelStep
	case tea.MouseButtonWheelDown:
		delta = wheelStep
	default:
		return
	}

	switch m.layout().paneAt(msg.X, msg.Y) {
	case listPane:
		m.scrollList(delta)
	case previewPane:
		m.preview.SetYOffset(m.preview.YOffset + delta)
	}
}

func (m *model) scrollList(delta int) {
	height := m.layout().listHeight
	m.listOffset = min(max(m.listOffset+delta, 0), max(len(m.rows)-height, 0))

	loc := m.current()
	cursor := min(max(loc.cursor, m.listOffset), m.listOffset+height-1)
	if cursor == loc.cursor || cursor >= len(m.rows) {
		return
	}
	loc.cursor = cursor
	m.showSelected()
}

func (m *model) handleKey(msg tea.KeyMsg) tea.Cmd {
	m.status = ""

	switch msg.String() {
	case "ctrl+c":
		return tea.Quit
	case "esc":
		return m.escape()
	case "up", "ctrl+p":
		m.moveCursor(-1)
		return nil
	case "down", "ctrl+n":
		m.moveCursor(1)
		return nil
	case "pgup":
		m.preview.PageUp()
		return nil
	case "pgdown":
		m.preview.PageDown()
		return nil
	case "tab":
		if m.suggestion != "" {
			return m.narrow()
		}
		return m.related()
	case "shift+tab":
		return m.back()
	case "enter":
		if m.suggestion != "" {
			return m.pivot()
		}
		return m.open()
	}

	before := m.input.Value()
	var cmd tea.Cmd
	m.input, cmd = m.input.Update(msg)
	if m.input.Value() == before {
		return cmd
	}

	return tea.Batch(cmd, m.edit(m.input.Value()))
}

func (m *model) current() *location {
	return &m.history[len(m.history)-1]
}

func (m *model) selected() (row, bool) {
	cursor := m.current().cursor
	if cursor < 0 || cursor >= len(m.rows) {
		return row{}, false
	}
	return m.rows[cursor], true
}

func (m *model) edit(query string) tea.Cmd {
	if m.current().note != nil {
		m.history = append(m.history, location{query: query})
		m.syncPlaceholder()
	} else {
		m.current().query = query
		m.current().cursor = 0
	}
	m.listOffset = 0
	m.syncFacets()

	return m.refresh(m.delay)
}

func (m *model) escape() tea.Cmd {
	if m.input.Value() != "" {
		m.input.SetValue("")
		return m.edit("")
	}
	if len(m.history) > 1 {
		return m.back()
	}
	return tea.Quit
}

func (m *model) related() tea.Cmd {
	selected, ok := m.selected()
	if !ok {
		return nil
	}
	m.history = append(m.history, location{note: selected.note})

	return m.jump()
}

func (m *model) narrow() tea.Cmd {
	base, _ := splitPartial(m.current().query)

	return m.follow(base, append(base[:len(base):len(base)], m.suggestion))
}

func (m *model) pivot() tea.Cmd {
	base, _ := splitPartial(m.current().query)

	return m.follow(base, []string{m.suggestion})
}

func (m *model) follow(base []string, terms []string) tea.Cmd {
	query := strings.Join(terms, " ") + " "
	m.current().cursor = 0
	if len(base) == 0 {
		m.current().query = query
	} else {
		m.current().query = strings.Join(base, " ")
		m.history = append(m.history, location{query: query})
	}

	return m.jump()
}

func (m *model) back() tea.Cmd {
	if len(m.history) <= 1 {
		return nil
	}
	m.history = m.history[:len(m.history)-1]

	return m.jump()
}

func (m *model) jump() tea.Cmd {
	m.rows = nil
	m.listOffset = 0
	m.syncInput()
	m.syncFacets()
	m.showSelected()

	return m.refresh(0)
}

func (m *model) open() tea.Cmd {
	selected, ok := m.selected()
	if !ok {
		return nil
	}
	command, err := editor.Command(selected.note.File, selected.line)
	if err != nil {
		m.status = err.Error()
		return nil
	}

	return tea.ExecProcess(command, func(err error) tea.Msg {
		return editorClosedMsg{err: err}
	})
}

func (m *model) refresh(delay time.Duration) tea.Cmd {
	m.generation++
	generation := m.generation

	return tea.Tick(delay, func(time.Time) tea.Msg {
		return searchDueMsg{generation: generation}
	})
}

func (m *model) search() tea.Cmd {
	m.searching = true
	book := m.book
	loc := m.target()
	generation := m.generation

	return func() tea.Msg {
		return searchDoneMsg{generation: generation, rows: lookup(book, loc)}
	}
}

func (m *model) moveCursor(delta int) {
	if len(m.rows) > 0 {
		loc := m.current()
		loc.cursor = ((loc.cursor+delta)%len(m.rows) + len(m.rows)) % len(m.rows)
	}
	m.clampCursor()
	m.showSelected()
}

func (m *model) clampCursor() {
	loc := m.current()
	if loc.cursor >= len(m.rows) {
		loc.cursor = len(m.rows) - 1
	}
	if loc.cursor < 0 {
		loc.cursor = 0
	}

	height := m.layout().listHeight
	if loc.cursor < m.listOffset {
		m.listOffset = loc.cursor
	}
	if loc.cursor >= m.listOffset+height {
		m.listOffset = loc.cursor - height + 1
	}
}

func (m *model) syncInput() {
	m.input.SetValue(m.current().query)
	m.input.CursorEnd()
	m.syncPlaceholder()
}

func (m *model) syncPlaceholder() {
	if note := m.current().note; note != nil {
		m.input.Placeholder = "Related to " + note.Title + " (type to search)"
		return
	}
	m.input.Placeholder = "Search notes, or #tag (showing recent notes)"
}

func (m *model) target() location {
	loc := *m.current()
	if m.suggestion != "" {
		base, _ := splitPartial(loc.query)
		loc.query = strings.Join(append(base[:len(base):len(base)], m.suggestion), " ")
	}
	return loc
}

func (m *model) indexTags() {
	m.tagTotals = tagCounts(allNotes(m.book))
	m.knownTags = make(map[string]bool, len(m.tagTotals))
	for tag := range m.tagTotals {
		m.knownTags[strings.ToLower(tag)] = true
	}
	m.baseCounts = nil
}

func (m *model) syncFacets() {
	m.facets, m.suggestion = m.currentFacets()

	query := m.current().query
	_, partial := splitPartial(query)
	if m.suggestion == "" || !hasPrefixFold(m.suggestion, partial) {
		m.input.SetSuggestions(nil)
		return
	}
	m.input.SetSuggestions([]string{query[:len(query)-len(partial)] + m.suggestion})
}

func (m *model) currentFacets() ([]facet, string) {
	loc := m.current()
	if loc.note != nil {
		counts := make(map[string]int, len(loc.note.Tags))
		for _, tag := range loc.note.Tags {
			counts[tag] = m.tagTotals[tag]
		}
		return rank(counts, nil), ""
	}

	base, partial := splitPartial(loc.query)
	if partial != "" && !m.knownTags[strings.ToLower(partial)] {
		candidates := complete(rank(m.countsWithin(base), base), partial, m.preferredTags(base))
		if len(candidates) > 0 {
			return candidates, candidates[0].tag
		}
	}

	return rank(tagCounts(rowNotes(m.rows)), strings.Fields(loc.query)), ""
}

func (m *model) countsWithin(base []string) map[string]int {
	if len(base) == 0 {
		return m.tagTotals
	}

	query := strings.Join(base, " ")
	if m.baseCounts == nil || m.baseQuery != query {
		m.baseQuery = query
		m.baseCounts = tagCounts(rowNotes(lookup(m.book, location{query: query})))
	}
	return m.baseCounts
}

func (m *model) preferredTags(base []string) []string {
	if len(base) > 0 || len(m.history) < 2 {
		return nil
	}
	if previous := m.history[len(m.history)-2]; previous.note != nil {
		return previous.note.Tags
	}
	return nil
}
