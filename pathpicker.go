package ntfilepicker

import (
	"io/fs"
	"path"
	"slices"
	"strings"

	"charm.land/bubbles/v2/key"
	tea "charm.land/bubbletea/v2"
	"github.com/charmbracelet/x/ansi"
)

// PathKeyMap is the keys a PathPicker takes. Other keys edit the path.
type PathKeyMap struct {
	Up, Down         key.Binding
	PageUp, PageDown key.Binding
	Top, Bottom      key.Binding
	Into, Parent     key.Binding
	Choose, Cancel   key.Binding
	Clear            key.Binding
}

// DefaultPathKeyMap is the compact path picker's default keys.
func DefaultPathKeyMap() PathKeyMap {
	k := DefaultKeyMap()
	return PathKeyMap{
		Up: k.Up, Down: k.Down, PageUp: k.PageUp, PageDown: k.PageDown,
		Top: k.Top, Bottom: k.Bottom,
		Into:   key.NewBinding(key.WithKeys("tab")),
		Parent: key.NewBinding(key.WithKeys("shift+tab")),
		Choose: key.NewBinding(key.WithKeys("enter")),
		Cancel: key.NewBinding(key.WithKeys("esc")),
		Clear:  key.NewBinding(key.WithKeys("ctrl+u")),
	}
}

// PathOption configures a PathPicker.
type PathOption func(*PathPicker)

// WithPathStyles sets the styles used by the prompt, input, rows and messages.
func WithPathStyles(s Styles) PathOption { return func(p *PathPicker) { p.core.styles = s } }

// WithPathKeyMap sets the path picker's keys.
func WithPathKeyMap(k PathKeyMap) PathOption { return func(p *PathPicker) { p.keys = k } }

// WithPathMarker draws a mark before each name, as WithMarker does.
func WithPathMarker(mark func(fs.DirEntry) string) PathOption {
	return func(p *PathPicker) { p.core.marker = mark }
}

// WithPathShowHidden shows dot-files even without a filter starting with a dot.
func WithPathShowHidden(show bool) PathOption {
	return func(p *PathPicker) { p.core.showHidden = show }
}

// WithPathDirsOnly lists and accepts directories alone.
func WithPathDirsOnly(on bool) PathOption { return func(p *PathPicker) { p.dirsOnly = on } }

// WithPathQuote quotes chosen paths for POSIX shells when needed (default true).
// Turn this off when passing the path directly to a filesystem or process API.
func WithPathQuote(on bool) PathOption { return func(p *PathPicker) { p.quote = on } }

// WithPathMaxRows limits the entry rows below the input (default 10).
func WithPathMaxRows(n int) PathOption { return func(p *PathPicker) { p.maxRows = max(1, n) } }

// WithPathPrompt sets the text before the path (default "> ").
func WithPathPrompt(s string) PathOption { return func(p *PathPicker) { p.core.prompt = s } }

// PathPicker completes a path in a small list, one directory at a time.
// Tab enters a directory; Enter accepts a file or directory; Esc cancels.
// The input is an unquoted path relative to the starting directory, or an
// absolute path from the filesystem's root. Chosen paths are always absolute.
// A host sets its size and places View wherever the popup belongs.
type PathPicker struct {
	core      Model // Shared listing, matching, input and safe name drawing.
	keys      PathKeyMap
	dirsOnly  bool
	quote     bool
	maxRows   int
	cancelled bool
}

// NewPathPicker makes a compact path picker over fsys, starting in dir.
func NewPathPicker(fsys fs.ReadDirFS, dir string, opts ...PathOption) PathPicker {
	p := PathPicker{
		core: New(fsys, dir, WithHome(""), WithPlaces(nil), WithPrompt("> ")),
		keys: DefaultPathKeyMap(), quote: true, maxRows: 10,
	}
	p.core.width, p.core.height = 40, 11
	for _, opt := range opts {
		opt(&p)
	}
	return p
}

// Init starts an asynchronous read of the current directory.
func (p PathPicker) Init() tea.Cmd { return p.core.load(p.core.viewDir()) }

// Chosen returns the accepted path, quoted if WithPathQuote is enabled.
func (p PathPicker) Chosen() (string, bool) { return p.core.Chosen() }

// ClearChosen forgets the choice after the host has taken it.
func (p *PathPicker) ClearChosen() { p.core.ClearChosen() }

// Cancelled reports whether Esc was pressed.
func (p PathPicker) Cancelled() bool { return p.cancelled }

// ClearCancelled forgets the cancellation after the host has handled it.
func (p *PathPicker) ClearCancelled() { p.cancelled = false }

// FilterValue returns the unquoted path being built, including directory parts.
func (p PathPicker) FilterValue() string { return p.core.FilterValue() }

// SetFilter sets an unquoted path and reads the directory it names.
func (p *PathPicker) SetFilter(s string) tea.Cmd {
	p.core.input.SetValue(s)
	p.core.input.CursorEnd()
	return p.edited()
}

func (p PathPicker) withFilter(s string) (PathPicker, tea.Cmd) {
	cmd := p.SetFilter(s)
	return p, cmd
}

func (p *PathPicker) edited() tea.Cmd {
	p.core.cursor, p.core.top = 0, 0
	p.core.moved = false
	p.rebuild()
	return p.core.load(p.core.viewDir())
}

// SetSize sets the maximum width and height, including the input row.
func (p *PathPicker) SetSize(w, h int) {
	p.core.width, p.core.height = max(1, w), max(1, h)
	p.reveal()
}

// Size returns the maximum room given to the widget.
func (p PathPicker) Size() (w, h int) { return p.core.Size() }

func (p PathPicker) bodyRows() int { return min(p.maxRows, p.core.height-1) }

func (p *PathPicker) reveal() {
	m := &p.core
	m.cursor = max(0, min(m.cursor, len(m.rows)-1))
	m.top = max(0, min(m.top, m.cursor))
	if h := p.bodyRows(); h > 0 && m.cursor >= m.top+h {
		m.top = m.cursor - h + 1
	}
}

func (p *PathPicker) rebuild() {
	_, rest := split(p.FilterValue())
	p.core.rows = p.core.rowsOf(p.core.viewDir(), rest, false)
	if p.dirsOnly {
		p.core.rows = slices.DeleteFunc(p.core.rows, func(r row) bool { return !r.entry.IsDir() })
	}
	// Keep fuzzy ranking within each group, with directories always first.
	slices.SortStableFunc(p.core.rows, func(a, b row) int {
		if a.entry.IsDir() == b.entry.IsDir() {
			return 0
		}
		if a.entry.IsDir() {
			return -1
		}
		return 1
	})
	p.reveal()
}

func (p *PathPicker) parent() tea.Cmd {
	dir := parentOf(p.core.viewDir())
	if dir != "/" {
		dir += "/"
	}
	return p.SetFilter(dir)
}

// Update handles directory replies, input and navigation. It never quits the
// host program; the host handles Chosen and Cancelled.
func (p PathPicker) Update(msg tea.Msg) (PathPicker, tea.Cmd) {
	m := &p.core
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		p.SetSize(msg.Width, msg.Height)
	case listedMsg:
		if msg.id == m.id && msg.epoch == m.epoch {
			m.listed(msg)
			p.rebuild()
		}
		return p, nil
	case tea.KeyPressMsg:
		_, rest := split(p.FilterValue())
		switch {
		case key.Matches(msg, p.keys.Cancel):
			p.cancelled = true
		case key.Matches(msg, p.keys.Choose):
			p.choose()
		case key.Matches(msg, p.keys.Into):
			if e := m.Current(); e != nil && e.IsDir() {
				prefix, _ := split(p.FilterValue())
				return p.withFilter(prefix + e.Name() + "/")
			}
		case key.Matches(msg, p.keys.Parent), msg.String() == "backspace" && rest == "":
			cmd := p.parent()
			return p, cmd
		case key.Matches(msg, p.keys.Clear):
			return p.withFilter("")
		case key.Matches(msg, p.keys.Up):
			m.cursor--
			m.moved = true
		case key.Matches(msg, p.keys.Down):
			m.cursor++
			m.moved = true
		case key.Matches(msg, p.keys.PageUp):
			m.cursor -= max(1, p.bodyRows())
			m.moved = true
		case key.Matches(msg, p.keys.PageDown):
			m.cursor += max(1, p.bodyRows())
			m.moved = true
		case key.Matches(msg, p.keys.Top):
			m.cursor = 0
			m.moved = true
		case key.Matches(msg, p.keys.Bottom):
			m.cursor = len(m.rows) - 1
			m.moved = true
		default:
			before := p.FilterValue()
			m.input, _ = m.input.Update(msg)
			if p.FilterValue() != before {
				cmd := p.edited()
				return p, cmd
			}
		}
		p.reveal()
	case tea.PasteMsg:
		m.input, _ = m.input.Update(msg)
		cmd := p.edited()
		return p, cmd
	}
	return p, nil
}

func (p *PathPicker) choose() {
	m := &p.core
	_, rest := split(p.FilterValue())
	entry := m.Current()
	if !m.moved {
		for _, r := range m.rows {
			if r.name == rest {
				entry = r.entry
				break
			}
		}
	}
	if entry == nil {
		return
	}
	m.chosen = path.Join(m.viewDir(), entry.Name())
	if p.quote {
		m.chosen = quotePath(m.chosen)
	}
}

func (p PathPicker) inputView() string {
	m := p.core
	value := []rune(p.FilterValue())
	pos := min(m.input.Position(), len(value))
	before, at, after := cleanName(string(value[:pos])), " ", ""
	if pos < len(value) {
		at, after = cleanName(string(value[pos])), cleanName(string(value[pos+1:]))
	}
	prompt := cleanName(m.prompt)
	room := max(0, m.width-ansi.StringWidth(prompt))
	if over := ansi.StringWidth(before+at) - room; over > 0 {
		before = ansi.TruncateLeft(before, over, "")
	}
	return m.styles.Prompt.Render(prompt) + m.styles.Input.Render(before) +
		m.styles.InputCursor.Render(at) + m.styles.Input.Render(after)
}

func quotePath(s string) string {
	if strings.IndexFunc(s, func(r rune) bool {
		return !(r >= 'a' && r <= 'z' || r >= 'A' && r <= 'Z' || r >= '0' && r <= '9' || strings.ContainsRune("/_-.", r))
	}) < 0 {
		return s
	}
	return "'" + strings.ReplaceAll(s, "'", "'\"'\"'") + "'"
}

// View draws the input followed by at most WithPathMaxRows entry rows.
func (p PathPicker) View() string {
	m := p.core
	lines := []string{p.inputView()}
	for i := m.top; i < len(m.rows) && len(lines) <= p.bodyRows(); i++ {
		r := m.rows[i]
		prefix, base := "  ", m.styles.File
		if r.entry.IsDir() {
			base = m.styles.Directory
		}
		if i == m.cursor {
			prefix, base = m.styles.Cursor.Render(">")+" ", m.styles.Selected
		}
		if m.marker != nil {
			prefix += m.styles.Mark.Render(cleanName(m.marker(r.entry))) + " "
		}
		lines = append(lines, prefix+m.name(r, m.width-ansi.StringWidth(prefix), base, m.styles.Match))
	}
	if len(m.rows) == 0 && p.bodyRows() > 0 {
		text, style := "  Empty folder.", m.styles.Empty
		l := m.lists[m.viewDir()]
		switch {
		case l == nil || !l.loaded:
			text = "  Reading…"
		case l.err != nil:
			text, style = "  ⚠ "+cleanName(problem(l.err)), m.styles.Error
		case p.FilterValue() != "":
			text = "  No matching files."
		}
		lines = append(lines, style.Render(text))
	}
	for i := range lines {
		lines[i] = ansi.Truncate(lines[i], m.width, "")
	}
	return strings.Join(lines, "\n")
}
