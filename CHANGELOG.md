# `ntfilepicker` CHANGELOG

## v0.2.0 (2026-10-28)

* Add `PathPicker`, a compact path completer with directory navigation, fuzzy
  filtering, shell quoting, configurable keys and styles, and an example in
  `examples/pathpick`. The existing `Model` is unchanged.
* Add a Path picker mode to the web demo, with action buttons, quoted results,
  and a Node smoke test over the pretend filesystem.
* Update Go dependencies in both modules, including go-booba v0.10.7 and
  Bubble Tea v2.1.1, and refresh the dependency notices.
* Fix gaps in the web demo's box-drawing borders through go-booba's updated
  ghostty-web renderer.

## v0.1.0 (2026-10-19)

First release.

* Extract gloss's file chooser (`internal/browse` at gloss e21ce3f) as its own
  module. The package is `ntfilepicker`; its test harness, script options and
  playground are internal (`internal/pickertest`, `internal/pickercfg`,
  `internal/cmd/play`) until the harness has a public API of its own.
* Add `WithUnchoosableNote`: the details of a file that cannot be chosen no
  longer say "gloss cannot show this file", but "This file cannot be chosen"
  or what the host sets.
* Add `examples/pick`, a complete chooser that prints the chosen path, and
  use it for the VHS recordings in place of gloss.
* Add a browser demo: `examples/demo` runs the chooser over a pretend home
  folder, natively and as WASM through booba, with a layout picker, action
  bar, and a Node smoke test. Examples are now a module of their own; the
  harness's in-memory filesystem moves to `internal/memfs` to share with it.
