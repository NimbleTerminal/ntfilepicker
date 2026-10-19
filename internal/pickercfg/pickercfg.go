// Package pickercfg is how the pickertest harness builds and reads the
// chooser: its options by name, and the state a script may check. It is
// shared by the tests and the playground.
package pickercfg

import (
	tea "charm.land/bubbletea/v2"
	"io/fs"
	"path"
	"slices"
	"strconv"
	"strings"

	"nimbleterminal.dev/ntfilepicker"
	"nimbleterminal.dev/ntfilepicker/internal/pickertest"
)

// Options a script may set, with `option NAME VALUE`:
//
//	layout list|columns|places   how folders are laid out
//	hidden true                  show dot-files
//	home /home/evan              where ~ leads (default /home/evan)
//	marker kinds                 mark folders and pictures with one-cell glyphs
//	marker emoji                 mark files by kind with emoji (two cells wide)
//	only .png                    only that extension may be chosen
//	columns N                    at most N columns
//	dirsonly true                Tab completes folders alone (as in go-to)
//	linksfollowed true           the host has made links to folders folders
//	filters kinds                offer the kinds Pictures (.png .jpg), Documents
//	                             (.md .txt), and Data (.csv); "active NAME,NAME"
//	                             starts with some chosen
const Options = "layout, hidden, home, marker, only, columns, filters, active, dirsonly, linksfollowed"

// StateKeys are the Probe values F2 records in the playground.
var StateKeys = []string{"dir", "filter", "current", "chosen", "layout", "types"}

// Config is the harness's way of making and reading the chooser.
func Config() pickertest.Config[ntfilepicker.Model] {
	return pickertest.Config[ntfilepicker.Model]{
		Strict: true,
		Message: func(msg tea.Msg) map[string]string {
			if v, ok := msg.(ntfilepicker.OpenWithMsg); ok {
				return map[string]string{"open-with": v.Path}
			}
			return nil
		},
		New: func(fsys *pickertest.FS, start string, o map[string]string) ntfilepicker.Model {
			opts := []ntfilepicker.Option{ntfilepicker.WithHome("/home/evan")}
			if v, ok := o["layout"]; ok {
				l, ok := ntfilepicker.ParseLayout(v)
				if !ok {
					panic("unknown layout " + v)
				}
				opts = append(opts, ntfilepicker.WithLayout(l))
			}
			if o["hidden"] == "true" {
				opts = append(opts, ntfilepicker.WithShowHidden(true))
			}
			if v, ok := o["home"]; ok {
				opts = append(opts, ntfilepicker.WithHome(v))
			}
			if o["marker"] == "kinds" {
				opts = append(opts, ntfilepicker.WithMarker(func(e fs.DirEntry) string {
					switch {
					case e.IsDir():
						return "▪"
					case strings.HasSuffix(e.Name(), ".png"):
						return "▫"
					}
					return " "
				}))
			}
			if o["marker"] == "emoji" {
				opts = append(opts, ntfilepicker.WithMarker(emojiMark))
			}
			if o["only"] != "" {
				ext := o["only"]
				opts = append(opts, ntfilepicker.WithSelectable(func(e fs.DirEntry) bool { return path.Ext(e.Name()) == ext }))
			}
			if o["linksfollowed"] == "true" {
				opts = append(opts, ntfilepicker.WithLinksFollowed())
			}
			if o["dirsonly"] == "true" {
				opts = append(opts, ntfilepicker.WithDirsOnlyCompletion(true))
			}
			if v, ok := o["columns"]; ok {
				n, _ := strconv.Atoi(v)
				opts = append(opts, ntfilepicker.WithMaxColumns(n))
			}
			if o["filters"] == "kinds" {
				kind := func(name, mark string, exts ...string) ntfilepicker.Filter {
					return ntfilepicker.Filter{Name: name, Mark: mark, Match: func(e fs.DirEntry) bool { return slices.Contains(exts, path.Ext(strings.ToLower(e.Name()))) }}
				}
				opts = append(opts, ntfilepicker.WithFilters([]ntfilepicker.Filter{
					kind("Pictures", "▫", ".png", ".jpg"), kind("Documents", "▪", ".md", ".txt"), kind("Data", "▤", ".csv"),
				}))
				if a := o["active"]; a != "" {
					opts = append(opts, ntfilepicker.WithActiveFilters(strings.Split(a, ",")...))
				}
			}
			return ntfilepicker.New(fsys, start, opts...)
		},
		Probe: func(m ntfilepicker.Model) map[string]string {
			chosen, _ := m.Chosen()
			current := ""
			if e := m.Current(); e != nil {
				current = e.Name()
			}
			return map[string]string{
				"dir":     m.Dir(),
				"viewdir": m.ViewDir(),
				"filter":  m.FilterValue(),
				"chosen":  chosen,
				"current": current,
				"rows":    strconv.Itoa(m.Rows()),
				"layout":  m.Layout().String(),
				"types":   strings.Join(m.ActiveFilters(), ","),
				"menu":    strconv.FormatBool(m.InMenu()),
				"sidebar": strconv.FormatBool(m.InSidebar()),
				"crumbs":  strings.Join(m.Crumbs(), " > "),
			}
		},
	}
}

// emojiMarks are sample marks of the kinds of file, by extension, two cells
// wide as emoji are; the chooser has none of its own, and a host brings its own.
var emojiMarks = map[string]string{
	".png": "📷", ".jpg": "📷", ".jpeg": "📷", ".gif": "📷", ".webp": "📷", ".heic": "📷",
	".svg": "🎨", ".mmd": "🧜", ".mermaid": "🧜", ".pdf": "📕", ".stl": "🧊", ".3mf": "🧊", ".obj": "🧊", ".ply": "🧊",
	".md": "📝", ".html": "📝", ".txt": "📝",
	".json": "🧾", ".ipynb": "📓", ".docx": "📄", ".xlsx": "📊", ".csv": "📊", ".grist": "📊",
}

// emojiMark is the mark of an entry: all of one width, two cells.
func emojiMark(e fs.DirEntry) string {
	if e.IsDir() {
		return "📁"
	}
	if m, ok := emojiMarks[strings.ToLower(path.Ext(e.Name()))]; ok {
		return m
	}
	return "  "
}
