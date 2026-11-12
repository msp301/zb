package browser

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/x/ansi"
)

const sideBySideWidth = 80
const keyHelp = "Up/Down select  PgUp/PgDn scroll  Tab related  Shift+Tab back  Enter edit  Esc clear/quit"

var (
	dimStyle      = lipgloss.NewStyle().Faint(true)
	accentStyle   = lipgloss.NewStyle().Foreground(lipgloss.Color("6"))
	tagStyle      = lipgloss.NewStyle().Foreground(lipgloss.Color("5"))
	titleStyle    = lipgloss.NewStyle().Bold(true)
	selectedStyle = lipgloss.NewStyle().Reverse(true)
	statusStyle   = lipgloss.NewStyle().Foreground(lipgloss.Color("3"))
)

type layout struct {
	sideBySide    bool
	listWidth     int
	listHeight    int
	previewWidth  int
	previewHeight int
}

func (m *model) layout() layout {
	body := max(m.height-3, 2)

	if m.width >= sideBySideWidth {
		listWidth := min(max(m.width/3, 30), 50)
		return layout{
			sideBySide:    true,
			listWidth:     listWidth,
			listHeight:    body,
			previewWidth:  m.width - listWidth - 3,
			previewHeight: body,
		}
	}

	listHeight := max(body/3, 1)
	return layout{
		listWidth:     max(m.width, 1),
		listHeight:    listHeight,
		previewWidth:  max(m.width, 1),
		previewHeight: max(body-listHeight-1, 1),
	}
}

func (m *model) resize() {
	l := m.layout()
	m.input.Width = max(m.width-4, 1)
	m.preview.Width = l.previewWidth
	m.preview.Height = l.previewHeight
	m.clampCursor()
	m.showSelected()
}

func (m *model) showSelected() {
	selected, ok := m.selected()
	if !ok {
		m.preview.SetContent("")
		return
	}

	width := m.layout().previewWidth
	var lines []string
	matchLine := 0
	if selected.note != nil {
		lines, matchLine = m.notePreview(selected, width)
	} else {
		lines = m.tagPreview(selected, width)
	}

	m.preview.SetContent(strings.Join(lines, "\n"))
	m.preview.SetYOffset(max(matchLine-m.preview.Height/3, 0))
}

func (m *model) notePreview(selected row, width int) ([]string, int) {
	note := selected.note
	header := fmt.Sprintf("%s:%d", m.relativePath(note.File), selected.line)
	lines := []string{accentStyle.Render(ansi.Truncate(header, width, "…"))}
	if len(note.Tags) > 0 {
		for _, tags := range wrap(strings.Join(note.Tags, " "), width) {
			lines = append(lines, tagStyle.Render(tags))
		}
	}
	lines = append(lines, "")

	matchIndex := selected.line - note.Start

	matchLine := 0
	for index, text := range strings.Split(note.Content, "\n") {
		gutter := "  "
		style := lipgloss.NewStyle()
		if index == 0 {
			style = titleStyle
		}
		if index == matchIndex && index > 0 {
			gutter = accentStyle.Render("│ ")
			style = titleStyle
			matchLine = len(lines)
		}

		for _, segment := range wrap(text, width-2) {
			lines = append(lines, gutter+style.Render(segment))
		}
	}

	return lines, matchLine
}

func (m *model) tagPreview(selected row, width int) []string {
	notes := taggedNotes(m.book, selected.tagId)
	lines := []string{
		tagStyle.Bold(true).Render(selected.tag),
		dimStyle.Render(count(len(notes), "note") + ", press Tab to browse"),
		"",
	}
	for _, note := range notes {
		lines = append(lines, ansi.Truncate(note.Title, width, "…"))
	}
	return lines
}

func (m *model) relativePath(file string) string {
	for _, dir := range m.options.Dirs {
		relative, err := filepath.Rel(dir, file)
		if err == nil && !strings.HasPrefix(relative, "..") {
			return relative
		}
	}
	return file
}

func wrap(text string, width int) []string {
	text = strings.ReplaceAll(text, "\t", "    ")
	if width < 1 {
		return []string{text}
	}
	return strings.Split(ansi.Wrap(text, width, ""), "\n")
}

func (m *model) View() string {
	if m.width == 0 || m.height == 0 {
		return ""
	}

	l := m.layout()
	lines := []string{m.input.View(), m.breadcrumb()}

	list := m.listLines(l)
	preview := strings.Split(m.preview.View(), "\n")
	if l.sideBySide {
		divider := dimStyle.Render(" │ ")
		for index := range list {
			line := list[index] + divider
			if index < len(preview) {
				line += preview[index]
			}
			lines = append(lines, line)
		}
	} else {
		lines = append(lines, list...)
		lines = append(lines, dimStyle.Render(strings.Repeat("─", l.listWidth)))
		lines = append(lines, preview...)
	}

	lines = append(lines, m.footer())

	return strings.Join(lines, "\n")
}

func (m *model) breadcrumb() string {
	labels := make([]string, 0, len(m.history))
	for _, loc := range m.history {
		labels = append(labels, loc.label())
	}
	trail := strings.Join(labels, " > ")
	if overflow := ansi.StringWidth(trail) - m.width; overflow > 0 {
		trail = ansi.TruncateLeft(trail, overflow+1, "…")
	}
	return dimStyle.Render(trail)
}

func (m *model) listLines(l layout) []string {
	lines := make([]string, 0, l.listHeight)
	cursor := m.current().cursor

	for index := m.listOffset; index < len(m.rows) && len(lines) < l.listHeight; index++ {
		lines = append(lines, renderRow(m.rows[index], index == cursor, l.listWidth))
	}

	if len(m.rows) == 0 {
		message := "No results"
		if m.applied != m.generation {
			message = "Searching…"
		}
		lines = append(lines, pad(dimStyle.Render(message), l.listWidth))
	}

	for len(lines) < l.listHeight {
		lines = append(lines, strings.Repeat(" ", l.listWidth))
	}

	return lines
}

func renderRow(r row, selected bool, width int) string {
	if selected {
		text := r.title
		if r.context != "" && r.context != r.title {
			text += "  " + r.context
		}
		return selectedStyle.Render(pad(ansi.Truncate(text, width, "…"), width))
	}

	if r.note == nil {
		return pad(tagStyle.Render(ansi.Truncate(r.title, width, "…")), width)
	}

	title := ansi.Truncate(r.title, width, "…")
	line := title
	remaining := width - ansi.StringWidth(title) - 2
	if r.context != "" && r.context != r.title && remaining > 3 {
		line += "  " + dimStyle.Render(ansi.Truncate(r.context, remaining, "…"))
	}
	return pad(line, width)
}

func pad(text string, width int) string {
	if missing := width - ansi.StringWidth(text); missing > 0 {
		return text + strings.Repeat(" ", missing)
	}
	return text
}

func (m *model) footer() string {
	if m.status != "" {
		return statusStyle.Render(ansi.Truncate(m.status, m.width, "…"))
	}

	return dimStyle.Render(ansi.Truncate(count(len(m.rows), "result")+"  |  "+keyHelp, m.width, "…"))
}

func count(number int, noun string) string {
	if number == 1 {
		return "1 " + noun
	}
	return fmt.Sprintf("%d %ss", number, noun)
}
