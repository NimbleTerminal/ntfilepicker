package main

import (
	"strings"
	"time"

	"charm.land/bubbles/v2/key"
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
	"github.com/charmbracelet/x/ansi"
	"nimbleterminal.dev/ntfilepicker"
	"nimbleterminal.dev/ntfilepicker/examples/internal/demorun"
	"nimbleterminal.dev/ntfilepicker/internal/memfs"
)

// pathDemo is a small host around the same pretend tree as the full browser.
// Keep the picker open after a choice so visitors can try another path.
type pathDemo struct {
	picker        ntfilepicker.PathPicker
	width, height int
	result        string
}

func newPathDemo() pathDemo {
	fsys := memfs.NewFS(tree()...).Latency(120 * time.Millisecond)
	m := pathDemo{
		picker: ntfilepicker.NewPathPicker(fsys, home),
		width:  80, height: 24,
	}
	m.resize()
	return m
}

func (m *pathDemo) resize() {
	m.picker.SetSize(min(58, max(1, m.width-2)), min(11, max(1, m.height-7)))
}

func (m pathDemo) Init() tea.Cmd { return m.picker.Init() }

func (m pathDemo) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width, m.height = max(1, msg.Width), max(1, msg.Height)
		m.resize()
		return m, nil
	case tea.KeyPressMsg:
		if msg.String() == "ctrl+c" {
			return m, tea.Quit
		}
	}
	var cmd tea.Cmd
	m.picker, cmd = m.picker.Update(msg)
	if chosen, ok := m.picker.Chosen(); ok {
		m.result = "Chosen: " + chosen
		m.picker.ClearChosen()
	}
	if m.picker.Cancelled() {
		m.result = "Cancelled. Type to try another path."
		m.picker.ClearCancelled()
		cmd = m.picker.SetFilter("")
	}
	return m, cmd
}

func (m pathDemo) View() tea.View {
	w, _ := m.picker.Size()
	popup := lipgloss.NewStyle().Width(w).Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("99")).Render(m.picker.View())
	result := m.result
	if result == "" {
		result = "Choose a path to see its shell quoting here."
	}
	lines := strings.Split("Complete a path\nStarting in "+home+" · pretend files only\n\n"+
		popup+"\n\n"+result, "\n")
	lines = lines[:min(len(lines), m.height)]
	for i := range lines {
		lines[i] = ansi.Truncate(lines[i], m.width, "")
	}
	v := tea.NewView(strings.Join(lines, "\n"))
	v.AltScreen = true
	return v
}

func runPathDemo() error {
	b := demorun.Binding
	return demorun.Run(newPathDemo(), "path", "Complete a path", demorun.Keys{
		Primary: []key.Binding{b("tab", "Into folder"), b("shift+tab", "Up a folder"), b("enter", "Choose path"), b("esc", "Cancel")},
		More: []key.Binding{b("up", "Previous"), b("down", "Next"), b("pgup", "Page up"), b("pgdown", "Page down"),
			b("ctrl+u", "Clear filter"), b("ctrl+c", "Quit")},
	})
}
