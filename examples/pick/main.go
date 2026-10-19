// Command pick chooses a file with ntfilepicker and prints its path.
//
//	go run ./examples/pick                 # start in the current folder
//	go run ./examples/pick ~/Downloads     # or in another
//	open "$(go run ./examples/pick)"       # the chooser draws on stderr
//
// Enter on a file chooses it; Esc clears the filter, or leaves without one.
package main

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"

	tea "charm.land/bubbletea/v2"
	"nimbleterminal.dev/ntfilepicker"
)

type model struct {
	picker ntfilepicker.Model
	volume string // The volume the paths are on, "" but on Windows.
	chosen string
}

func (m model) Init() tea.Cmd { return m.picker.Init() }

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyPressMsg:
		switch msg.String() {
		case "ctrl+c":
			return m, tea.Quit
		case "esc":
			if m.picker.FilterValue() != "" {
				return m, m.picker.ClearFilter()
			}
			return m, tea.Quit
		}
	}
	var cmd tea.Cmd
	m.picker, cmd = m.picker.Update(msg)
	if path, ok := m.picker.Chosen(); ok {
		m.chosen = m.volume + filepath.FromSlash(path)
		return m, tea.Quit
	}
	return m, cmd
}

func (m model) View() tea.View {
	v := tea.NewView(m.picker.View())
	v.AltScreen = true
	v.MouseMode = tea.MouseModeCellMotion
	return v
}

func main() {
	dir := "."
	if len(os.Args) > 1 {
		dir = os.Args[1]
	}
	abs, err := filepath.Abs(dir)
	if err != nil {
		fmt.Fprintln(os.Stderr, "pick:", err)
		os.Exit(1)
	}
	// The chooser's paths are slash-separated from the filesystem's root.
	volume := filepath.VolumeName(abs)
	root := os.DirFS(volume + string(filepath.Separator)).(fs.ReadDirFS)
	m := model{picker: ntfilepicker.New(root, filepath.ToSlash(abs[len(volume):])), volume: volume}

	final, err := tea.NewProgram(m, tea.WithOutput(os.Stderr)).Run()
	if err != nil {
		fmt.Fprintln(os.Stderr, "pick:", err)
		os.Exit(1)
	}
	if chosen := final.(model).chosen; chosen != "" {
		fmt.Println(chosen)
		return
	}
	os.Exit(2)
}
