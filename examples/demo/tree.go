package main

import (
	"fmt"
	"strconv"
	"time"
)

// home is where the demo starts: a pretend person's home folder.
const home = "/home/ada"

// tree is the pretend filesystem, in the lines memfs takes: a path, a size,
// and a date. Nothing in it is real, and nothing has contents.
func tree() []string {
	lines := []string{
		file("home/ada/Desktop/todo.txt", 312, "2026-10-18"),
		file("home/ada/Desktop/Screenshot 2026-10-17 at 09.41.12.png", 812004, "2026-10-17"),

		file("home/ada/Documents/résumé.pdf", 148220, "2026-09-03"),
		file("home/ada/Documents/budget 2026.xlsx", 40210, "2026-10-11"),
		file("home/ada/Documents/notes.md", 6120, "2026-10-19"),
		file("home/ada/Documents/日本語のメモ.txt", 2048, "2026-06-21"),
		file("home/ada/Documents/Letters/to the landlord.docx", 18800, "2026-02-14"),
		file("home/ada/Documents/Letters/thank you, grandma.docx", 15100, "2025-12-28"),
		file("home/ada/Documents/Recipes/🍝 carbonara.md", 1900, "2026-03-08"),
		file("home/ada/Documents/Recipes/🥐 croissants (attempt 4).md", 3400, "2026-05-30"),
		file("home/ada/Documents/Recipes/🌮 tacos al pastor.md", 2600, "2026-08-02"),
		file("home/ada/Documents/Taxes/2024.pdf", 91000, "2025-04-02"),
		file("home/ada/Documents/Taxes/2025.pdf", 88000, "2026-04-01"),
		folder("home/ada/Documents/Taxes/receipts/"),

		file("home/ada/Pictures/logo.svg", 4200, "2026-05-05"),
		file("home/ada/Pictures/Wallpapers/aurora.heic", 3900000, "2026-01-09"),
		file("home/ada/Pictures/Wallpapers/dunes.heic", 4100000, "2026-01-09"),
		file("home/ada/Pictures/Wallpapers/tide pools.heic", 3700000, "2026-01-10"),

		file("home/ada/Music/The Quiet Modems/Dial Tone/01 Handshake.flac", 31000000, "2025-11-01"),
		file("home/ada/Music/The Quiet Modems/Dial Tone/02 Carrier Lost.flac", 28000000, "2025-11-01"),
		file("home/ada/Music/The Quiet Modems/Dial Tone/03 56k Lullaby.flac", 35000000, "2025-11-01"),
		file("home/ada/Music/Ferrite Hearts/Core Memory/01 Magnetic.mp3", 7200000, "2024-07-19"),
		file("home/ada/Music/Ferrite Hearts/Core Memory/02 Bit Rot.mp3", 6800000, "2024-07-19"),
		file("home/ada/Music/Sølvklang/Nordlys/01 Polarnatt.flac", 40100000, "2026-02-02"),

		file("home/ada/Videos/holiday 2026.mov", 1873000000, "2026-08-20"),
		file("home/ada/Videos/cat discovers the printer.mp4", 48000000, "2026-04-11"),

		file("home/ada/Downloads/setup-tool-3.2.dmg", 52000000, "2026-09-30"),
		file("home/ada/Downloads/invoice (final) v2 (1).pdf", 31000, "2026-09-29"),
		file("home/ada/Downloads/invoice (final) v2.pdf", 31000, "2026-09-29"),
		file("home/ada/Downloads/dataset.csv", 12400000, "2026-10-02"),
		file("home/ada/Downloads/slides.pdf", 2200000, "2026-10-14"),

		file("home/ada/projects/website/index.html", 5400, "2026-10-12"),
		file("home/ada/projects/website/style.css", 2100, "2026-10-12"),
		file("home/ada/projects/website/images/hero.png", 640000, "2026-10-12"),
		file("home/ada/projects/picker/README.md", 3100, "2026-10-19"),
		file("home/ada/projects/picker/go.mod", 412, "2026-10-19"),
		file("home/ada/projects/picker/main.go", 2600, "2026-10-19"),
		file("home/ada/projects/picker/main_test.go", 1800, "2026-10-19"),
		file("home/ada/projects/picker/testdata/scripts/basic.txt", 400, "2026-10-19"),
		file("home/ada/projects/a folder whose name is far too long to fit in any narrow column/x.txt", 1, "2026-01-01"),

		file("home/ada/.profile", 120, "2025-01-01"),
		file("home/ada/.config/editor/settings.json", 900, "2026-07-07"),
		link("home/ada/notes", "/home/ada/Documents/notes.md"),
		link("home/ada/recipes", "/home/ada/Documents/Recipes"),

		file("etc/hosts", 200, "2025-06-01"),
		file("usr/local/bin/pick", 4200000, "2026-10-19"),
	}
	// A camera roll to scroll through.
	day := time.Date(2026, time.June, 1, 0, 0, 0, 0, time.UTC)
	for i := 1; i <= 240; i++ {
		date := day.AddDate(0, 0, i/3).Format("2006-01-02")
		lines = append(lines, file(fmt.Sprintf("home/ada/Pictures/Camera/IMG_%04d.jpg", 4000+i), 1500000+int64(i*7919%900000), date))
	}
	// And a folder of thousands, to filter.
	for i := range 2000 {
		date := time.Date(2021, time.January, 1, 0, 0, 0, 0, time.UTC).AddDate(0, 0, i).Format("2006-01-02")
		lines = append(lines, file("home/ada/Archive/logs/app-"+date+".log", 2000+int64(i*37%50000), date))
	}
	return lines
}

// file, folder and link write one line each, quoting the path as memfs wants
// a path with spaces in it.
func file(path string, size int64, date string) string {
	return fmt.Sprintf("%s %d @%s", strconv.Quote(path), size, date)
}

func folder(path string) string { return strconv.Quote(path) }

func link(path, target string) string { return strconv.Quote(path) + " -> " + target }
