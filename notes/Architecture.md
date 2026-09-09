---
title: Architecture
tags:
  - codebase
  - architecture
aliases:
  - Node design
---

# Architecture

Part of [[Codebase Overview]].

## Node shape

One sealed interface carries the whole tree:

```go
type Node interface {
	renderNode(*renderer)
	validateNode(*validation, string)
}
```

Three concrete cases live in `node.go`/`script_code.go`/`unsafe_html.go`:
`textNode` (escaped text), `rawNode` (unescaped, `UnsafeHTML`), and
`elementNode` (tag + props + children). `Node` is sealed by the unexported
`renderNode`/`validateNode` methods — only files in `package html` can add a
new case.

> [!warning] This blocks a naive package split
> Every element constructor returns `elementNode{...}` directly — an
> unexported struct. Splitting the ~130 generated element files into a
> subpackage would require exporting `elementNode`'s fields, a public-API
> change, not a navigation cleanup. See [[Decisions Log#Why nothing here got split]].

## Render / validate pipeline

```mermaid
graph LR
    A[html.Div/P/... builders] --> B[elementNode tree]
    B --> C[Validate walks tree]
    C --> D[emitWarning per issue]
    B --> E[Render walks tree]
    E --> F[escaped HTML string]
```

- `Render(n)` always runs `Validate(n)` first and emits every warning via
  the package-level `WarningHandler` (default: `log.Printf`) before writing
  any output — warnings never block rendering.
- `elementNode.validateNode` checks void-element-with-children, duplicate
  attributes, and duplicate `id` values in one tree walk.
- `renderer.text`/`raw` are the only two places bytes reach the output
  buffer un-mediated by tag structure — `text` HTML-escapes, `raw` doesn't
  (that's the whole contract behind `UnsafeHTML`).

## Attribute assembly

Every generated constructor funnels through one shared function,
`globalAttrs` (`common_props.go`) — the busiest node in the graph
(degree 137, see [[Decisions Log#graphify]]). It applies `id`/`class`/
`style`/`title` plus `customAttributes`, which regex-filters caller-supplied
`Attribute`s against a valid HTML attribute-name pattern before they reach
the tree.

## Elements

Every `<tag>.go` file is hand-written. `cmd/htmlgen` used to write these from
an `index.json` element schema, but that schema never existed in this repo —
the generator was dead code from the first commit and was removed (see
[[Decisions Log#htmlgen removal]]). Every element with real HTML attributes
exposes them as typed fields; elements with no attributes beyond the generic
four remain unchanged. `script`
is special-cased to `ScriptTag`/`script_tag.go` to avoid colliding with
`ScriptCode` (`script_code.go`), the hand-written raw-JS-body node.
