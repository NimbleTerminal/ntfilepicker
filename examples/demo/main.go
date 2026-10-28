// Command demo is ntfilepicker over a pretend home folder, in a terminal or
// in the browser. Nothing it lists is real: the tree is made up in tree.go,
// and folders take a moment to read, as a slow disk's would, to show that the
// chooser does not wait on them.
//
//	go run ./demo                  # or -layout columns, -layout list
package main

import (
	"flag"
	"fmt"
	"io/fs"
	"os"
	"path"
	"strings"
	"time"

	"charm.land/bubbles/v2/key"
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
	"github.com/dustin/go-humanize"
	"nimbleterminal.dev/ntfilepicker"
	"nimbleterminal.dev/ntfilepicker/examples/internal/demorun"
	"nimbleterminal.dev/ntfilepicker/internal/memfs"
)

// below is the rows under the chooser: a rule and the last file picked. They
// are below so that the chooser's mouse rows are the screen's.
const below = 2

var (
	ruleStyle   = lipgloss.NewStyle().Foreground(lipgloss.Color("238"))
	hintStyle   = lipgloss.NewStyle().Foreground(lipgloss.Color("244"))
	pickedStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("212")).Bold(true)
)

type model struct {
	fsys   *memfs.FS
	picker ntfilepicker.Model
	width  int
	picked string // The last file chosen, and what is known of it.
}

func newModel(layout ntfilepicker.Layout) model {
	fsys := memfs.NewFS(tree()...).Latency(120 * time.Millisecond)
	picker := ntfilepicker.New(fsys, home,
		ntfilepicker.WithHome(home),
		ntfilepicker.WithLayout(layout),
		ntfilepicker.WithMarker(mark),
		ntfilepicker.WithFilters(filters),
		ntfilepicker.WithSelectable(func(e fs.DirEntry) bool { return path.Ext(e.Name()) != ".dmg" }),
		ntfilepicker.WithUnchoosableNote("Installers cannot be picked in this demo"),
		ntfilepicker.WithPrompt("  Pick a file: "),
	)
	return model{fsys: fsys, picker: picker, width: 80}
}

func (m model) Init() tea.Cmd { return m.picker.Init() }

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		msg.Height = max(1, msg.Height-below)
		var cmd tea.Cmd
		m.picker, cmd = m.picker.Update(msg)
		return m, cmd
	case tea.KeyPressMsg:
		switch msg.String() {
		case "ctrl+c":
			return m, tea.Quit
		case "esc":
			// In a browser there is nowhere to leave to: Esc only clears.
			return m, m.picker.ClearFilter()
		}
	}
	var cmd tea.Cmd
	m.picker, cmd = m.picker.Update(msg)
	if chosen, ok := m.picker.Chosen(); ok {
		m.picker.ClearChosen()
		m.picked = m.describe(chosen)
	}
	return m, cmd
}

// describe says what was picked, as a host would act on it.
func (m model) describe(p string) string {
	info, err := m.fsys.Stat(strings.TrimPrefix(p, "/"))
	if err != nil {
		return p
	}
	return fmt.Sprintf("%s  ·  %s  ·  %s", p, humanize.Bytes(uint64(max(info.Size(), 0))), info.ModTime().Format("2 Jan 2006"))
}

func (m model) View() tea.View {
	status := hintStyle.Render("Enter picks a file. Nothing here is real: a pretend home folder, read slowly on purpose.")
	if m.picked != "" {
		status = hintStyle.Render("Picked ") + pickedStyle.Render(m.picked)
	}
	rule := ruleStyle.Render(strings.Repeat("─", max(m.width, 1)))
	v := tea.NewView(m.picker.View() + "\n" + rule + "\n" + status)
	v.AltScreen = true
	v.MouseMode = tea.MouseModeCellMotion
	return v
}

// mark is a one-glance kind for each entry, as a host such as gloss gives.
func mark(e fs.DirEntry) string {
	if e.IsDir() {
		return "📁"
	}
	if e.Type()&fs.ModeSymlink != 0 {
		return "🔗"
	}
	for _, k := range kinds {
		for _, ext := range k.exts {
			if strings.EqualFold(path.Ext(e.Name()), ext) {
				return k.mark
			}
		}
	}
	return "📄"
}

type kind struct {
	name, mark string
	exts       []string
}

var kinds = []kind{
	{"Pictures", "📷", []string{".jpg", ".png", ".heic", ".svg"}},
	{"Documents", "📝", []string{".pdf", ".md", ".txt", ".docx", ".xlsx", ".csv"}},
	{"Music", "🎵", []string{".flac", ".mp3"}},
	{"Videos", "🎬", []string{".mov", ".mp4"}},
	{"Code", "🧩", []string{".go", ".mod", ".html", ".css", ".json"}},
}

var filters = func() []ntfilepicker.Filter {
	var out []ntfilepicker.Filter
	for _, k := range kinds {
		out = append(out, ntfilepicker.Filter{Name: k.name, Mark: k.mark, Match: func(e fs.DirEntry) bool { return mark(e) == k.mark }})
	}
	return out
}()

func main() {
	name := flag.String("layout", "places", "the layout to start in: places, columns, or list")
	widget := flag.String("widget", "browser", "the widget to try: browser or path")
	flag.Parse()
	if *widget == "path" {
		if err := runPathDemo(); err != nil {
			fmt.Fprintln(os.Stderr, "demo:", err)
			os.Exit(1)
		}
		return
	}
	if *widget != "browser" {
		fmt.Fprintf(os.Stderr, "demo: no widget %q\n", *widget)
		os.Exit(2)
	}
	layout, ok := ntfilepicker.ParseLayout(*name)
	if !ok {
		fmt.Fprintf(os.Stderr, "demo: no layout %q\n", *name)
		os.Exit(2)
	}
	b := demorun.Binding
	if err := demorun.Run(newModel(layout), "browse", "Pick a file", demorun.Keys{
		Primary: []key.Binding{b("ctrl+l", "Layout"), b("ctrl+g", "Places"), b("ctrl+f", "File types"), b("alt+up", "Up a folder")},
		More:    []key.Binding{b("alt+left", "Back"), b("alt+right", "Forward"), b("tab", "Complete"), b("esc", "Clear filter"), b("ctrl+c", "Quit")},
	}); err != nil {
		fmt.Fprintln(os.Stderr, "demo:", err)
		os.Exit(1)
	}
}
