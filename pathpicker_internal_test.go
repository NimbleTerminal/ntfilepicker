package ntfilepicker

import (
	"io/fs"
	"testing"
	"testing/fstest"
)

func TestPathPickerLateListings(t *testing.T) {
	fsys := fstest.MapFS{"one/a": &fstest.MapFile{}, "two/b": &fstest.MapFile{}}
	p := NewPathPicker(fsys, "/one")
	first := p.Init()
	second := p.SetFilter("/two/")
	p, _ = p.Update(second())
	p, _ = p.Update(first())
	if len(p.core.rows) != 1 || p.core.rows[0].name != "b" {
		t.Fatal("late reply changed the current directory")
	}
	other := NewPathPicker(fstest.MapFS{"two/wrong": &fstest.MapFile{}}, "/two")
	p, _ = p.Update(other.Init()())
	p, _ = p.Update(listedMsg{id: p.core.id, epoch: p.core.epoch - 1, dir: "/two", err: fs.ErrPermission})
	if p.core.lists["/two"].err != nil || p.core.rows[0].name != "b" {
		t.Fatal("accepted another picker or epoch's reply")
	}
	if cmd := p.SetFilter("/one/"); cmd != nil {
		t.Fatal("cached directory was read again")
	}
}
