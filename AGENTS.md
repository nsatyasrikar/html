# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Commands

```sh
go build ./...                          # build everything
go test ./...                           # run all tests
go test -run TestRenderTree ./...       # run a single test by name
go vet ./...                            # static check
```

No lint config, no external dependencies, no network/API key needed to build or test.

## Architecture

Single `package html` at repo root. Every HTML element is a hand-written
pair: a `<Tag>Props` struct (`ID`, `Class`, `Style`, `Title`,
`Attributes []Attribute`, plus any element-specific fields) and a
constructor function (`Div`, `P`, `Input`, ...) that builds an `elementNode`
and returns it as the sealed `Node` interface (`node.go`). Most element
files only expose the generic four props; a handful (`input.go`,
`button.go`, `a.go`, `meta.go`, `link.go`, ...) also expose element-specific
fields where a consumer needed them. **Rule: extend an element file on
demand when a consumer needs a specific attribute — never pre-extend
elements speculatively.** Reach anything not yet exposed through the loose
`Attributes` escape hatch until that day comes. `render.go`, `validate.go`, and the
other pipeline files are hand-written for the same reason: they're core
pipeline, not an element.

Rendering flow: `Render` (`render.go`) always runs `Validate` (`validate.go`)
first, emits every `Warning` through the package-level `WarningHandler`
(`warning.go`), then walks the tree writing escaped HTML. `elementNode`,
`textNode` (`node.go`), and `rawNode`/`UnsafeHTML` (`unsafe_html.go`) are the
only three `Node` implementations — the interface is sealed by unexported
methods, so a new node kind can only be added from inside `package html`.
Every constructor funnels its common props through `globalAttrs`
(`common_props.go`), which also filters custom `Attribute`s through a
name-validity regex (`attribute.go`) before they reach the tree.

See `notes/Architecture.md` and `notes/Decisions Log.md` for the reasoning
behind these boundaries (in particular why element files were never split
into a subpackage).

## graphify

This project has a knowledge graph at graphify-out/ with god nodes, community structure, and cross-file relationships.

Rules:
- For codebase questions, first run `graphify query "<question>"` when graphify-out/graph.json exists. Use `graphify path "<A>" "<B>"` for relationships and `graphify explain "<concept>"` for focused concepts. These return a scoped subgraph, usually much smaller than GRAPH_REPORT.md or raw grep output.
- If graphify-out/wiki/index.md exists, use it for broad navigation instead of raw source browsing.
- Read graphify-out/GRAPH_REPORT.md only for broad architecture review or when query/path/explain do not surface enough context.
- After modifying code, run `graphify update .` to keep the graph current (AST-only, no API cost).


## obsidian

Codebase notes live in `notes/` as an Obsidian vault (see `notes/Codebase Overview.md` for the index).

Rules:
- When adding a new note, link it into `notes/Codebase Overview.md`'s map of notes so it's reachable from the index.
- Before writing a new note, check existing notes for overlapping topics and `[[wikilink]]` to them instead of duplicating content.
- When a change touches something an existing note documents (e.g. a decision, a design constraint), update that note in the same commit rather than leaving it stale.
- Prefer linking related notes to each other directly (not just through the index) when they reference the same concept, so navigation works from any entry point.
