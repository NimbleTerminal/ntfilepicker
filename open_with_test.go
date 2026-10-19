package ntfilepicker_test

import (
	tea "charm.land/bubbletea/v2"
	"nimbleterminal.dev/ntfilepicker"
	"nimbleterminal.dev/ntfilepicker/internal/pickertest"
	"testing"
)

func TestOpenWithIntent(t *testing.T) {
	for _, layout := range []ntfilepicker.Layout{ntfilepicker.LayoutList, ntfilepicker.LayoutColumns, ntfilepicker.LayoutPlaces} {
		d := pickertest.New(t, ntfilepicker.New(pickertest.NewFS("d/a.txt"), "/d", ntfilepicker.WithLayout(layout)), 80, 20)
		d.Press("alt+enter")
		found := false
		for _, msg := range d.Messages {
			if v, ok := msg.(ntfilepicker.OpenWithMsg); ok {
				found = true
				if v.Path != "/d/a.txt" || v.Directory {
					t.Fatalf("%+v", v)
				}
			}
		}
		if !found {
			t.Fatal("no intent", layout)
		}
		if _, ok := d.M.Chosen(); ok {
			t.Fatal("selected file")
		}
		_, _ = d.M.Update(tea.KeyPressMsg{Code: 'w', Text: "w"})
	}
}
