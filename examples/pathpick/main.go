// Command pathpick completes a path in a compact popup and prints it.
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
	picker ntfilepicker.PathPicker
	volume string
	chosen string
}

func (m model) Init() tea.Cmd { return m.picker.Init() }

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.picker.SetSize(min(60, msg.Width), min(11, msg.Height))
		return m, nil
	case tea.KeyPressMsg:
		if msg.String() == "ctrl+c" {
			return m, tea.Quit
		}
	}
	var cmd tea.Cmd
	m.picker, cmd = m.picker.Update(msg)
	if chosen, ok := m.picker.Chosen(); ok {
		m.chosen = m.volume + filepath.FromSlash(chosen)
		return m, tea.Quit
	}
	if m.picker.Cancelled() {
		return m, tea.Quit
	}
	return m, cmd
}

func (m model) View() tea.View { return tea.NewView(m.picker.View()) }

func main() {
	dir := "."
	if len(os.Args) > 1 {
		dir = os.Args[1]
	}
	abs, err := filepath.Abs(dir)
	if err != nil {
		fmt.Fprintln(os.Stderr, "pathpick:", err)
		os.Exit(1)
	}
	volume := filepath.VolumeName(abs)
	root := os.DirFS(volume + string(filepath.Separator)).(fs.ReadDirFS)
	// Print the raw native path; quoting is for a host inserting shell text.
	m := model{
		picker: ntfilepicker.NewPathPicker(root, filepath.ToSlash(abs[len(volume):]), ntfilepicker.WithPathQuote(false)),
		volume: volume,
	}
	final, err := tea.NewProgram(m, tea.WithOutput(os.Stderr)).Run()
	if err != nil {
		fmt.Fprintln(os.Stderr, "pathpick:", err)
		os.Exit(1)
	}
	if chosen := final.(model).chosen; chosen != "" {
		fmt.Println(chosen)
		return
	}
	os.Exit(2)
}
