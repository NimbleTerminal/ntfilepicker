package ntfilepicker_test

import (
	"fmt"
	"testing"

	"nimbleterminal.dev/ntfilepicker"
	"nimbleterminal.dev/ntfilepicker/internal/pickertest"
)

// huge is a folder of n files beside one with a few, as a photo dump is.
func huge(n int) *pickertest.FS {
	f := pickertest.NewFS("data/small/a.txt")
	for i := 0; i < n; i++ {
		f.Add(fmt.Sprintf("data/big/IMG_%06d.jpg 1000", i))
	}
	return f
}

func benchmark(b *testing.B, layout ntfilepicker.Layout, keys string) {
	f := huge(100000)
	d := pickertest.New(b, ntfilepicker.New(f, "/data/big", ntfilepicker.WithLayout(layout)), 120, 40)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		d.Type(keys)
		_ = d.Screen()
		d.Press("ctrl+u")
	}
}

func BenchmarkTypeInHugeFolderList(b *testing.B) { benchmark(b, ntfilepicker.LayoutList, "img_00123") }
func BenchmarkTypeInHugeFolderColumns(b *testing.B) {
	benchmark(b, ntfilepicker.LayoutColumns, "img_00123")
}

func BenchmarkViewHugeFolderColumns(b *testing.B) {
	f := huge(100000)
	d := pickertest.New(b, ntfilepicker.New(f, "/data", ntfilepicker.WithLayout(ntfilepicker.LayoutColumns)), 120, 40)
	// The cursor starts on big/, whose listing is the last column.
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = d.M.View()
	}
}
