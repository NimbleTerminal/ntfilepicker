package ntfilepicker_test

import (
	"testing"

	"nimbleterminal.dev/ntfilepicker"
	"nimbleterminal.dev/ntfilepicker/internal/pickertest"
)

// quiet leaves out the places, which the chooser would describe too.
var quiet = []ntfilepicker.Option{ntfilepicker.WithHome("/h"), ntfilepicker.WithPlaces(nil)}

func linksFS() *pickertest.FS {
	return pickertest.NewFS("d/real/", "d/a -> /d/real", "d/b -> /d/real", "d/f -> /d/file", "d/file 1", "d/gone -> /nowhere")
}

// A link is looked at once per listing, with Stat, which does not open it; and
// not at all by a chooser whose host has already done so.
func TestLinksAreDescribedOnceAndNeverOpened(t *testing.T) {
	f := linksFS()
	d := pickertest.New(t, ntfilepicker.New(f, "/d", quiet...), 60, 8)
	if got := f.Stats(); len(got) != 4 { // a, b, f, gone.
		t.Errorf("the chooser described %d links, want 4: %v", len(got), got)
	}
	if got := f.Opens(); len(got) != 0 {
		t.Errorf("the chooser opened %v", got)
	}
	if !d.Screen().Contains("a/") || !d.Screen().Contains("b/") {
		t.Errorf("links to a folder are folders:\n%s", d.Screen())
	}

	f = linksFS()
	d = pickertest.New(t, ntfilepicker.New(f, "/d", append(quiet, ntfilepicker.WithLinksFollowed())...), 60, 8)
	if got := f.Stats(); len(got) != 0 {
		t.Errorf("a chooser told links are followed described %v", got)
	}
	if got := f.Opens(); len(got) != 0 {
		t.Errorf("the chooser opened %v", got)
	}
	_ = d
}

// FollowLinks, which a host that filters listings runs first, is the same: it
// describes, and opens nothing.
func TestFollowLinksDoesNotOpen(t *testing.T) {
	f := linksFS()
	entries, _ := f.ReadDir("d")
	out := ntfilepicker.FollowLinks(f, "/d", entries)
	folders := 0
	for _, e := range out {
		if e.IsDir() {
			folders++
		}
	}
	if folders != 3 { // real, a, b.
		t.Errorf("%d folders, want 3", folders)
	}
	if len(f.Opens()) != 0 || len(f.Stats()) != 4 {
		t.Errorf("opens %v, stats %v", f.Opens(), f.Stats())
	}
}
