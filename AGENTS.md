# ntfilepicker

A file chooser component for Bubble Tea v2, extracted from gloss. Module:
`nimbleterminal.dev/ntfilepicker`, Go 1.26.8+. README is the introduction,
DEVELOP.md is for working on it; keep them, this file and CHANGELOG consistent.

## Layout

- The root package `ntfilepicker`: the chooser. Layouts (`list`, `columns`,
  `places`) only draw one state; keys, filter and completion are the same in each.
  `PathPicker` is a separate compact path completer, sharing listing, matching,
  caching and safe name drawing with `Model`; its keys and view are its own.
- `internal/pickertest`: the script-driven harness (in-memory FS, key driver,
  golden screens, VHS tapes, play/record/replay). Generic beyond its FS; see its
  README's "Taking it out" before making it public.
- `internal/pickercfg`: how a script's `option` lines build the chooser, and what
  `state` reads. Its emoji marks are a sample, not any host's.
- `internal/memfs`: the in-memory filesystem the harness and the demo share.
- `internal/cmd/play`: the playground.
- `examples` (its own module): `pick`, a complete host over the real
  filesystem; `pathpick`, a compact path picker; `demo`, the chooser over a
  pretend home folder, native and WASM through booba; `internal/demorun`, the
  booba action-bar wrapper.
- `web`: the demo page (`index.html`, `site.css`, `site.js` tracked; the rest
  generated). `wasm.work`: Bubble Tea/clipboard forks for WASM only.
- `scripts`: the site check and the Node smoke test of the WASM demo.
- `internal/tools/notices`: `THIRD_PARTY_NOTICES.md` for both builds.
- `testdata/scripts`: the scripts and their `.golden` screens.
  `testdata/pathpicker` holds the compact widget's scripts and screens.

## Conventions

- The chooser knows nothing of any host. What a host adds (marks, sort, which
  files may be chosen, wording about them) comes in by options; Esc is the host's.
  `PathPicker` handles Esc by setting `Cancelled`, which the host acts on.
- Any `fs.ReadDirFS`; paths are slash-separated from its root. Reads are
  commands, and stale replies are dropped by id and epoch.
- Every screen must stay sound: no row wider than the terminal by either width
  table, no control characters, and nothing from a file name reaching the
  terminal as an escape (`Screen.Problems`, the sweep test).
- Change behavior by adding a script first; the golden diff is the review.
- Upstream Bubble Tea only, no `replace` directives: hosts bring their forks.
- Follow the gloss code style: gofmt, small public API, comments in plain words.

## Verification

`task ci` before finishing (GOWORK=off): fmt, tidy/verify of both modules, race
tests, vet, Windows builds, the examples, the WASM site with its checks and smoke
test (Node 18+), notices, govulncheck (network). `task test` and `task screens`
while working. The demo's tree is made up: keep real names and paths out of it.
