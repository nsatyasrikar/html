---
title: Decisions Log
tags:
  - codebase
  - decisions
---

# Decisions Log

Part of [[Codebase Overview]]. Non-obvious choices, recorded so nobody
re-derives them from scratch later.

## Why nothing here got split

`client` (sibling repo) split its 1034 generated marker types out of the
core AST package because they were auxiliary — implementation detail behind
a sealed `Type` interface, reachable only via an exported `TypeMarker`
escape hatch. This repo's ~130 generated files are the opposite: each one
*is* the public API (`html.Div`, `html.P`, ...), called directly by every
consumer. Moving them to a subpackage would turn `html.Div(...)` into
`elements.Div(...)` everywhere — an API break, not a cleanup. See
[[Architecture#Node shape]]. No split; `cmd/htmlgen` staying a separate
`main` package (it only runs at generation time, never imported by `html`)
is the only boundary that already made sense.

## Publishing as a public module

Made go-gettable: added `LICENSE` (MIT), a `README.md` with install
instructions and runnable examples (every example verified against actual
`go run`/`go test` output, not assumed — see [[Codebase Overview#Quick facts]]),
and tagged/pushed `v0.1.0` so `go get github.com/nsatyasrikar/html` resolves.

`go.mod`'s `go` directive was `1.26.4` — an exact patch pin, which forces a
toolchain download for any consumer on a different patch version. Loosened
to `go 1.21`, the real floor (nothing here needs newer than generics,
1.18+). Public library go.mod directives should stay a minor version, not a
patch, unless a specific toolchain fix is required.

README's Limitations section is the honest list of what's *not* generated
yet (`href` on `a`, `name`/`content` on `meta`, etc. — see
[[Architecture#Generation]]) — keep it in sync with `cmd/htmlgen` rather
than letting it go stale.

## graphify

A knowledge graph of this codebase lives at `graphify-out/` (gitignored —
[[Codebase Overview]]), rebuilt automatically by a post-commit git hook
(`graphify hook install`, already wired here). AST-only extraction, no
LLM/API cost for a code-only corpus like this one. Query it with
`graphify query "<question>"` before grepping raw source for architecture
questions.
