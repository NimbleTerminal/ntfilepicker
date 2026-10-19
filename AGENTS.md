# ntfilepicker

A file chooser component for Bubble Tea v2, extracted from gloss. Module:
`nimbleterminal.dev/ntfilepicker`, Go 1.26.8+. README is the introduction,
DEVELOP.md is for working on it; keep them, this file and CHANGELOG consistent.

## Layout

- The root package `ntfilepicker`: the chooser. Layouts (`list`, `columns`,
  `places`) only draw one state; keys, filter and completion are the same in each.
- `internal/pickertest`: the script-driven harness (in-memory FS, key driver,
  golden screens, VHS tapes, play/record/replay). Generic beyond its FS; see its
  README's "Taking it out" before making it public.
- `internal/pickercfg`: how a script's `option` lines build the chooser, and what
  `state` reads. Its emoji marks are a sample, not any host's.
- `internal/cmd/play`: the playground. `examples/pick`: a complete host.
- `testdata/scripts`: the scripts and their `.golden` screens.

## Conventions

- The chooser knows nothing of any host. What a host adds (marks, sort, which
  files may be chosen, wording about them) comes in by options; Esc is the host's.
- Any `fs.ReadDirFS`; paths are slash-separated from its root. Reads are
  commands, and stale replies are dropped by id and epoch.
- Every screen must stay sound: no row wider than the terminal by either width
  table, no control characters, and nothing from a file name reaching the
  terminal as an escape (`Screen.Problems`, the sweep test).
- Change behavior by adding a script first; the golden diff is the review.
- Upstream Bubble Tea only, no `replace` directives: hosts bring their forks.
- Follow the gloss code style: gofmt, small public API, comments in plain words.

## Verification

`task ci` before finishing (GOWORK=off): fmt, tidy/verify, race tests, vet,
Windows build, the example, govulncheck (network). `task test` and
`task screens` while working.
