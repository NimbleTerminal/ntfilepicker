package ntfilepicker_test

import (
	"io/fs"
	"strconv"
	"testing"
	"testing/fstest"

	"charm.land/bubbles/v2/key"
	tea "charm.land/bubbletea/v2"
	"nimbleterminal.dev/ntfilepicker"
	"nimbleterminal.dev/ntfilepicker/internal/pickertest"
)

func TestPathPickerScripts(t *testing.T) {
	pickertest.RunDir(t, pickertest.Config[ntfilepicker.PathPicker]{
		Strict: true,
		New: func(fsys *pickertest.FS, start string, o map[string]string) ntfilepicker.PathPicker {
			return ntfilepicker.NewPathPicker(fsys, start,
				ntfilepicker.WithPathDirsOnly(o["dirsonly"] == "true"),
				ntfilepicker.WithPathShowHidden(o["hidden"] == "true"),
				ntfilepicker.WithPathQuote(o["quote"] != "false"))
		},
		Probe: func(p ntfilepicker.PathPicker) map[string]string {
			chosen, _ := p.Chosen()
			return map[string]string{"filter": p.FilterValue(), "chosen": chosen, "cancelled": strconv.FormatBool(p.Cancelled())}
		},
	}, "testdata/pathpicker")
}

func TestPathPickerOptionsAndChoice(t *testing.T) {
	fsys := fstest.MapFS{
		"folder/file": &fstest.MapFile{},
		"a'b $x":      &fstest.MapFile{},
		".hidden":     &fstest.MapFile{},
	}
	for _, tc := range []struct {
		name         string
		opts         []ntfilepicker.PathOption
		filter, want string
	}{
		{"quoted", nil, "a'b", "'/a'\"'\"'b $x'"},
		{"raw", []ntfilepicker.PathOption{ntfilepicker.WithPathQuote(false)}, "a'b", "/a'b $x"},
		{"directory", []ntfilepicker.PathOption{ntfilepicker.WithPathDirsOnly(true)}, "", "/folder"},
		{"hidden", []ntfilepicker.PathOption{ntfilepicker.WithPathShowHidden(true)}, "hid", "/.hidden"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			d := pickertest.New(t, ntfilepicker.NewPathPicker(fsys, "/", tc.opts...), 30, 5)
			d.Run(d.M.SetFilter(tc.filter))
			d.Press("enter")
			if got, ok := d.M.Chosen(); !ok || got != tc.want {
				t.Fatalf("choice = %q, %v; want %q", got, ok, tc.want)
			}
			d.M.ClearChosen()
			if _, ok := d.M.Chosen(); ok {
				t.Fatal("choice was not cleared")
			}
			d.Press("esc")
			if !d.M.Cancelled() {
				t.Fatal("Esc did not cancel")
			}
			d.M.ClearCancelled()
			if d.M.Cancelled() {
				t.Fatal("cancellation was not cleared")
			}
		})
	}
}

func TestPathPickerSafeScreens(t *testing.T) {
	fsys := fstest.MapFS{
		"\x1b[31mred\a":      &fstest.MapFile{},
		"家族👨‍👩‍👧‍👦":          &fstest.MapFile{},
		"❤️longlonglonglong": &fstest.MapFile{},
	}
	for w := 1; w <= 42; w++ {
		for h := 1; h <= 6; h++ {
			d := pickertest.New(t, ntfilepicker.NewPathPicker(fsys, "/",
				ntfilepicker.WithPathMarker(func(fs.DirEntry) string { return "📷\a" })), w, h)
			d.StrictFit = true
			d.Screen()
			d.Run(d.M.SetFilter("\x1b[31m"))
			d.Screen()
			d.Run(d.M.SetFilter("❤️👨‍👩‍👧‍👦\n"))
			for range []rune(d.M.FilterValue()) {
				d.Press("left")
				d.Screen()
			}
			if gotW, gotH := d.M.Size(); gotW != w || gotH != h {
				t.Fatalf("size = %dx%d, want %dx%d", gotW, gotH, w, h)
			}
		}
	}
}

func TestPathPickerEditingAndKeys(t *testing.T) {
	fsys := fstest.MapFS{"a": &fstest.MapFile{}, "b": &fstest.MapFile{}, "c": &fstest.MapFile{}}
	k := ntfilepicker.DefaultPathKeyMap()
	k.Choose = key.NewBinding(key.WithKeys("alt+enter"))
	d := pickertest.New(t, ntfilepicker.NewPathPicker(fsys, "/",
		ntfilepicker.WithPathKeyMap(k), ntfilepicker.WithPathMaxRows(1), ntfilepicker.WithPathPrompt("Path: ")), 20, 8)
	d.StrictFit = true
	d.Press("pgdown", "enter")
	if _, ok := d.M.Chosen(); ok {
		t.Fatal("Enter bypassed the custom key map")
	}
	d.Press("alt+enter")
	if got, _ := d.M.Chosen(); got != "/b" {
		t.Fatalf("page down chose %q, want /b", got)
	}
	if got := d.Screen().Lines(); len(got) > 2 {
		t.Fatalf("max rows exceeded: %q", got)
	}
	d.Press("pgup", "ctrl+u")
	d.Type("ac")
	d.Press("left")
	d.Type("b")
	if got := d.M.FilterValue(); got != "abc" {
		t.Fatalf("editing in the middle = %q", got)
	}
	d.Press("ctrl+u")
	d.Send(tea.PasteMsg{Content: "b"})
	d.Press("alt+enter")
	if got, _ := d.M.Chosen(); got != "/b" {
		t.Fatalf("paste chose %q", got)
	}
	d.M.ClearChosen()
	d.Run(d.M.SetFilter("no match"))
	d.Press("alt+enter")
	if _, ok := d.M.Chosen(); ok {
		t.Fatal("chose a nonexistent path")
	}
}
