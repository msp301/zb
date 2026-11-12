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
	if strings.TrimSpace(loc.query) == "" {
		return "recent"
	}
	return loc.query
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

	_, err = tea.NewProgram(newModel(options, book), tea.WithAltScreen()).Run()
	return err
}

func newModel(options Options, book *notebook.Notebook) *model {
	input := textinput.New()
	input.Prompt = "> "
	input.Cursor.SetMode(cursor.CursorStatic)
	input.SetValue(options.Query)
	input.CursorEnd()
	input.Focus()

	m := &model{
		options: options,
		book:    book,
		input:   input,
		preview: viewport.New(0, 0),
		history: []location{{query: options.Query}},
		delay:   searchDelay,
	}
	m.rows = lookup(book, m.history[0])
	m.syncInput()

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
		return m, m.refresh(0)

	case tea.KeyMsg:
		return m, m.handleKey(msg)
	}

	return m, nil
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
		return m.pivot()
	case "shift+tab":
		return m.back()
	case "enter":
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

func (m *model) pivot() tea.Cmd {
	selected, ok := m.selected()
	if !ok {
		return nil
	}

	if selected.note != nil {
		m.history = append(m.history, location{note: selected.note})
	} else {
		m.history = append(m.history, location{query: selected.tag})
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
	m.showSelected()

	return m.refresh(0)
}

func (m *model) open() tea.Cmd {
	selected, ok := m.selected()
	if !ok {
		return nil
	}
	if selected.note == nil {
		return m.pivot()
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
	loc := *m.current()
	generation := m.generation

	return func() tea.Msg {
		return searchDoneMsg{generation: generation, rows: lookup(book, loc)}
	}
}

func (m *model) moveCursor(delta int) {
	m.current().cursor += delta
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
