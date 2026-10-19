# Developing ntfilepicker

The chooser is a component of its own: it was written for gloss and knows nothing
of it, and takes what a host adds (marks, sort, which files may be chosen) as options, over
any `fs.ReadDirFS`. In gloss, `internal/app/opener.go` is the adapter. Folders are read
by commands and cached; a *layout* only draws that state (`list`, `columns`,
`places`), so keys, filter, and completion are the same in each. Paths are
slash-separated from the filesystem's root.

Its tests drive it through `internal/pickertest` (see its README), which works for any
component with `Init`, `Update`, and `View`: an in-memory filesystem (with
latency, injected read errors, and read counts), a driver that sends keys and
clicks and settles the commands that follow, and scripts. A script is a file in
`testdata/scripts`: an `fs` part listing files and a `script`
part of commands (`press`, `type`, `click`, `snapshot`, `state`, ...; see
`pickertest.RunDir`). Snapshots are compared with the `.golden` file beside the
script, which is a readable screen:

```sh
task test       # run them
task screens    # write the golden screens
```

Scripts understand the words of [VHS](https://github.com/charmbracelet/vhs)
tapes where the two overlap (`Type "re"`, `Down 2`, `Ctrl+L`, `Alt+Up`,
`Screenshot`; recording commands such as `Sleep` and `Set` are ignored), so a
tape reads as one. What VHS has no word for is ours: `fs`, `state`, `expect`,
`reads`, `click`. A script can also become a tape, to record a GIF of the
behavior it tests:

```sh
task tape SCRIPT=columns     # dist/columns.tape
task gif SCRIPT=columns      # records dist/columns.gif (needs vhs, ttyd, ffmpeg)
task gifs                    # the showcase scripts
```

The tasks are the usual ones: `task test`, `task screens` (rewrite
the golden screens), and `task bench`. `CMD`, `KEYS`, and `PAUSE` change
what the tape runs, the keys it substitutes, and its pace; the tool behind
them is `go run ./internal/pickertest/tape` (see its `-h`).

The tape types the command, presses the keys, and leaves out what a recording
cannot do (checks, clicks, resizes, the `fs` part: it runs in a real folder).
VHS has no Home or End, and takes Alt only with a character, so a script's
`alt+up` is left as a comment unless `-keys alt+up=Ctrl+Up,alt+left=Ctrl+O`
stands in keys the browser also answers to. With `vhs` installed, a test
checks that every script makes a tape VHS accepts.

### Playing with it

The harness can also be driven by a person, on the same in-memory trees the
scripts use:

```sh
task play                                        # the chooser over a demo tree
task play -- -from ~/Downloads                   # over a copy of a real folder (names, sizes, dates)
task play -- -script testdata/scripts/places.txt # a script's tree and options
task record NAME=thing                           # play, and keep what you do as a script
task replay SCRIPT=places                        # watch a script step by step, its checks shown
```

While recording, `F1` takes a snapshot, `F2` records checks of the state
(`state dir …`, `state current …`), `F3` leaves a note to edit; `Ctrl-C` ends and
writes the script, and `task record` then writes its golden screens. The
recording keeps the typing, keys, clicks, and wheel as the script words for them,
at a fixed size (90x20) so goldens stay small and alike, and puts the tree it
used in the script's `fs` part (look before sharing a `-from` copy: names and
sizes are real). In a replay `Space` does a step and then its checks, `p` plays
by itself (`+`/`-` change the pace), `r` starts again, `q` leaves; the status line
shows each check as it passes or fails, and the command exits non-zero if any did.

### Soundness checks

Every settled screen of a script is checked (`Screen.Problems`) for the faults
that mess up a terminal: more rows than there are; a row wider than the terminal
by **either** of two width tables, grapheme clusters and wcwidth, which disagree
about some emoji (a family of people joined with zero-width joiners is two cells
to one and six to the other), so a row that fits by one and overflows by the
other wraps on some terminals; and a control character in the text, as a file
name can carry. `col TEXT N` asserts the cell column something is drawn at, which
a wide character shifts from where letters would say, and `reject-raw TEXT`
asserts that nothing a name says reaches the terminal as an escape sequence.
`TestNoSizeOrNameMakesAnUnsoundScreen` sweeps every layout over 26 widths and 9
heights with awkward names (CJK, emoji, joined and selected sequences, combining
marks, very long names, bells, escapes) and fails on any unsound screen. Names
are drawn as names (`names.go`): control bytes as their symbols (`␇`, `␛`), and a
cluster the two tables count differently as its first character.

To try an idea, add a script, run it with `-update-screens`, and read the
golden file; then keep it. After changing how anything is drawn, the diff of
the golden files is the review. `task bench` times typing and drawing in a folder of 100,000 files.

## Examples and the browser demo

`examples` is a module of its own, so go-booba and the demo's needs stay out
of the library's `go.mod`; it points at the library with a `replace`.
`examples/pick` is a chooser over the real filesystem; `examples/demo` is the
chooser over a pretend home folder (`tree.go`, built with `internal/memfs`,
whose reads are slowed on purpose), run by booba natively and in the browser.

```sh
task demo                     # in this terminal (-- -layout columns)
task build-wasm-site          # web/: the WASM demo, booba's runtime, licenses
task serve-wasm-site          # http://localhost:8000/ntfilepicker/
task check-site               # the site's assets, and the demo driven under Node
```

The browser build uses `wasm.work`, which adds Bubble Tea and clipboard forks
that build for js/wasm; native builds use upstream. `web/index.html`,
`site.css` and `site.js` are tracked; the rest of `web/` is generated. The
Pages workflow runs `task ci` and publishes `web/`.
`scripts/smoke-demo.mjs` drives the real WASM binary through the booba bridge:
it filters, opens a folder and picks a file. `task notices` rewrites
`THIRD_PARTY_NOTICES.md` from both builds' dependencies.

## Verification

`task ci` is the gate: formatting, `go mod tidy`/`verify` of both modules, the
race tests, vet, Windows builds, the examples, the WASM site and its checks,
the notices, and `govulncheck` of native and WASM builds (needs the network).
It runs with `GOWORK=off`, so a local workspace cannot hide a dependency problem.

The library uses upstream Bubble Tea. gloss builds it with its own renderer fork
and in its browser demo; the scripts' golden screens do not depend on either.
