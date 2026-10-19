# ntfilepicker

<p>
    <a href="https://github.com/NimbleTerminal/ntfilepicker/tags"><img src="https://img.shields.io/github/tag/NimbleTerminal/ntfilepicker.svg" alt="Latest Release"></a>
    <a href="https://pkg.go.dev/nimbleterminal.dev/ntfilepicker"><img src="https://pkg.go.dev/badge/nimbleterminal.dev/ntfilepicker.svg" alt="GoDoc"></a>
</p>

A file chooser for [Bubble Tea](https://github.com/charmbracelet/bubbletea) v2:
a folder shown as a list, as Finder-style columns, or as columns with a sidebar
of places; filtered by what you type; walked with the arrow keys, Tab completion,
breadcrumbs, or the mouse. It is the chooser behind
[gloss](https://github.com/NimbleTerminal/gloss)'s `o` key, and knows nothing of
gloss: marks for the kinds of file, the sort, and which files may be chosen
come in by options. Any `fs.ReadDirFS` will do.

Module `nimbleterminal.dev/ntfilepicker`, Go 1.26.8+.

```sh
go get nimbleterminal.dev/ntfilepicker
go run nimbleterminal.dev/ntfilepicker/examples/pick@latest ~/Downloads
```

## Use

A host drives it like any bubble, and owns Esc:

```go
picker := ntfilepicker.New(os.DirFS("/").(fs.ReadDirFS), "/home/me",
	ntfilepicker.WithLayout(ntfilepicker.LayoutColumns),
	ntfilepicker.WithSelectable(func(e fs.DirEntry) bool { return strings.HasSuffix(e.Name(), ".pdf") }),
)

// In Update:
if key, ok := msg.(tea.KeyPressMsg); ok && key.String() == "esc" {
	if picker.FilterValue() != "" {
		return m, picker.ClearFilter()
	}
	return m, tea.Quit
}
picker, cmd = picker.Update(msg)
if path, ok := picker.Chosen(); ok {
	// path is slash-separated from the filesystem's root.
}
```

[`examples/pick`](examples/pick/main.go) is a complete program: it prints the
chosen path, and draws on stderr so that `$(pick)` works.

Folders are read by commands and cached; reads that come back after the
chooser has moved on are dropped. Names are drawn as names: control bytes as
their symbols, and nothing a file name says reaches the terminal as an escape.

## Options

| Option | Does |
| --- | --- |
| `WithLayout` | `LayoutList`, `LayoutColumns`, or `LayoutPlaces` (Ctrl-L cycles) |
| `WithMarker` | a mark before each name, such as an emoji for its kind |
| `WithSort` | the order of entries |
| `WithSelectable` | which files may be chosen; others are dimmed |
| `WithUnchoosableNote` | what the details of a file that cannot be chosen say |
| `WithFilters`, `WithActiveFilters` | a menu of kinds of file to show alone (Ctrl-F) |
| `WithPlaces` | the sidebar's places, in place of `DefaultPlaces` |
| `WithShowHidden`, `WithHome`, `WithPrompt`, `WithFooter`, `WithMaxColumns` | as named |
| `WithDirsOnlyCompletion` | Tab completes folders alone |
| `WithLinksFollowed` | the host already presents links to folders as folders (`FollowLinks`) |
| `WithStyles`, `WithKeyMap` | colors and keys (`DefaultStyles`, `DefaultKeyMap`) |

Alt-Enter sends an `OpenWithMsg` for the highlighted entry without choosing it,
for a host that can hand files to other apps.

## Keys

| Key | Does |
| --- | --- |
| typing | filters the folder |
| ↑ ↓, PgUp PgDn, Home End | move |
| Enter | into a folder, or choose a file |
| → | into a folder, or accept the completion shown |
| ← or Backspace (empty filter), Alt-↑ | up a folder |
| Tab, Shift-Tab | complete |
| Alt-← Alt-→ | back and forward through folders visited |
| Ctrl-L | next layout |
| Ctrl-G | to the sidebar of places, and back |
| Ctrl-F | the menu of kinds of file |

## Developing

See [DEVELOP.md](DEVELOP.md): the tests are scripts with readable golden
screens, which also play in the terminal and record as VHS GIFs. `task ci`
is the gate.

## License

MIT, see [LICENSE.txt](LICENSE.txt).
