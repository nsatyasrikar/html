---
title: Codebase Overview
tags:
  - codebase
  - index
aliases:
  - html
  - Overview
---

# html — Codebase Overview

Index note for the `html` Go module. Start here.

> [!abstract] What this is
> A typed Go builder for HTML: one exported function per element
> (`html.Div`, `html.P`, `html.Input`, ...) returning a `Node`, rendered to a
> string via `Render`. Every element file is hand-written — see
> [[Architecture#Elements]].

## Map of notes

- [[Architecture]] — Node/render/validate pipeline, element-file pattern, package boundary
- [[Decisions Log]] — non-obvious choices and why, for future-you

## Package layout at a glance

| Package | Path | Role |
|---|---|---|
| `html` | repo root | Node interface, renderer, validator, ~140 hand-written element constructors |

## Quick facts

- Module: `github.com/nsatyasrikar/html`, published under MIT, tagged `v0.1.0` — see [[Decisions Log#Publishing as a public module]]
- Requires Go 1.21+, no third-party dependencies
- No API key or network access needed to build/test
- Every element file (`div.go`, `p.go`, ...) follows one template: a `<Name>Props` struct plus a constructor calling `globalAttrs` + `elementNode{...}` — all hand-written, see [[Decisions Log#Closing the Limitations list]]
- Knowledge graph lives at `graphify-out/` (gitignored, regenerates via `graphify update .` — see [[Decisions Log#graphify]])

## Related code

- [README.md](../README.md) — install instructions and runnable examples
- [LICENSE](../LICENSE)
- [go.mod](../go.mod)
