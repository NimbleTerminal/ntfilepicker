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

[Try it in your browser](https://nimbleterminal.github.io/ntfilepicker/), over a
pretend home folder: nothing is read from your computer.

Module `nimbleterminal.dev/ntfilepicker`, Go 1.26.8+.

```sh
go get nimbleterminal.dev/ntfilepicker
```

To try the examples from a clone:

```sh
go -C examples run ./pick ~/Downloads   # prints the path you pick
go -C examples run ./pathpick           # compact path completion
go -C examples run ./demo               # the browser demo, in your terminal
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
[`examples/demo`](examples/demo/main.go) is the browser demo: marks, kinds of
file, a file that cannot be chosen, and status lines below the chooser.

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

## Compact path picker

`PathPicker` lists one directory in a small popup. The host controls its size
and placement; it has no sidebar, layouts, or history.

```go
picker := ntfilepicker.NewPathPicker(os.DirFS("/").(fs.ReadDirFS), "/home/me")
picker.SetSize(48, 8)
cmd := picker.Init() // Run this command, and any command returned by Update.

// In Update:
picker, cmd = picker.Update(msg)
if path, ok := picker.Chosen(); ok {
	// Insert path into shell input, then call picker.ClearChosen().
}
if picker.Cancelled() {
	// Close the popup, then call picker.ClearCancelled() if reusing it.
}
// In View: place picker.View() where the popup belongs.
```

Typing fuzzy-filters names, with directories first. Tab enters the highlighted
directory and completes that part of `FilterValue`; Shift-Tab goes up.
Backspace goes up when the final path component is empty, including after `/`.
↑/↓, PgUp/PgDn and Home/End move; Enter chooses a file or directory; Esc cancels.
An exact name takes priority until you move the highlight yourself.
Ctrl-U clears the whole input; Ctrl-W deletes a word.

The filter is an unquoted path, relative to the starting directory or beginning
with `/` for the filesystem's root. Choices are slash-separated absolute paths
from that root. Spaces and shell specials are quoted with POSIX single quotes,
including embedded apostrophes. Use `WithPathQuote(false)` for raw paths passed
to filesystem or process APIs, or when the host supplies its own quoting.
[`examples/pathpick`](examples/pathpick/main.go) prints a raw native path.
Try the [PathPicker web demo](https://nimbleterminal.github.io/ntfilepicker/?widget=path)
over the pretend home folder, or run `go -C examples run ./demo -widget path`.

Options are `WithPathStyles(Styles)`, `WithPathKeyMap(PathKeyMap)`,
`WithPathMarker`, `WithPathShowHidden`, `WithPathDirsOnly`, `WithPathQuote`,
`WithPathMaxRows` (10 entry rows by default), and `WithPathPrompt` (`"> "`).
`SetFilter` returns a command; always run it. Directory reads are asynchronous
and cached, and names use the same safe rendering as the full chooser.

## Developing

See [DEVELOP.md](DEVELOP.md): the tests are scripts with readable golden
screens, which also play in the terminal and record as VHS GIFs. `task demo`
runs the demo; `task serve-wasm-site` serves the browser one. `task ci` is the gate.

## License

Released under the [MIT License](./LICENSE.txt).  Part of the Nimble Terminal platform.

Copyright (c) 2026 [Neomantra Corp](https://www.neomantra.com).

----
Made with :heart: and :fire: by the team behind [Nimble Terminal](https://nimbleterminal.com).
