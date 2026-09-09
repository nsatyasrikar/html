# Close README Limitations

Date: 2026-09-09

## Background

README's Limitations section (see `notes/Decisions Log.md#htmlgen removal`)
currently lists three gaps:

1. Most element files only expose the generic four props (`ID`, `Class`,
   `Style`, `Title`) plus `Attributes`; element-specific attributes (`href`
   on `a`, `name`/`content` on `meta`, ...) are hand-added on demand, so
   most of ~135 element files don't expose their real attributes yet.
2. "No component abstraction beyond `Node`" — composition is plain Go
   functions returning `Node`, plus `Fragment(children...) []Node`.
3. No CSS/JS scoping — `<style>`/`<script>` content is emitted as-is.

Goal: close all three enough to delete them from Limitations, without
reversing the just-made decision to kill `cmd/htmlgen` (see
`notes/Decisions Log.md#htmlgen removal`) — no schema, no generator tool
gets reintroduced.

## Scope decisions (from brainstorming)

- **Attribute coverage**: hand-write every element's real attributes
  directly into its file — same pattern as `input.go`/`a.go`/`meta.go`/
  `link.go` today (optional `*string`/`*bool` fields, conditional
  `append(a, Attribute{...})` in the constructor). No schema file, no
  `cmd/` tool. A one-off local script may be used during implementation to
  move faster across ~130 near-identical edits, but nothing generator-shaped
  is committed — every file remains independently hand-editable and the
  script is discarded, not checked in.
- **Component abstraction**: no code change. `Fragment` (added in commit
  `d0a70b6`) plus plain functions returning `Node`/`[]Node` already are the
  component model for this package. This bullet is dropped from Limitations
  as already resolved.
- **CSS/JS scoping**: add an opt-in `ScopedStyle` helper that namespaces a
  CSS string's selectors under one scope id. JS scoping is not attempted —
  reworded to a narrower, honestly-scoped caveat (see below) rather than
  claimed solved.

## Design

### 1. Element-specific attributes

For each element file, add fields for that tag's real HTML5 attributes
(per MDN's per-element attribute tables), following the existing template:

```go
type FooProps struct {
    ID, Class, Style, Title, Bar, Baz *string
    Qux                                *bool
    Attributes                         []Attribute
}

func Foo(p FooProps, children ...Node) Node {
    a := globalAttrs(p.ID, p.Class, p.Style, p.Title, p.Attributes)
    if p.Bar != nil {
        a = append(a, Attribute{"bar", *p.Bar, false})
    }
    if p.Qux != nil {
        a = append(a, Attribute{"qux", "true", true})
    }
    return elementNode{name: "foo", props: a, void: false, children: children}
}
```

Boolean attributes follow `input.go`'s existing convention: `Attribute{name,
"true", true}` (the third field marks it boolean so `render.go` emits the
bare attribute name, no value).

Elements with genuinely no attributes beyond the global four (most inline
text-level elements — `b.go`, `i.go`, `em.go`, `strong.go`, `mark.go`,
`small.go`, ...) are left untouched; adding fields nothing uses would
violate the same YAGNI call this project keeps making. "Closing" limitation
#1 means: every element that *has* real HTML attributes exposes them: it
does not mean giving every file matching fields it doesn't need.

Deprecated/obsolete elements already present in this repo (`acronym.go`,
`big.go`, `center.go`, `dir.go`, `font.go`, `frame.go`, `frameset.go`,
`marquee.go`, `nobr.go`, `noframes.go`, `plaintext.go`, `strike.go`,
`tt.go`, `xmp.go`, `rb.go`, `rtc.go`) get their legacy-but-real attributes
(e.g. `bgcolor` on `body.go`/`table.go`, `size`/`color` on `font.go`) since
those are still what makes them meaningfully different from a `div` — same
treatment as any other element, not special-cased.

Implementation note: a throwaway script (in the session scratchpad, never
committed) may batch the mechanical part of this — inserting fields and
`if` blocks from a lookup table of tag → attribute list — but every
resulting file must build and read exactly as if it had been typed by
hand, matching the existing files' style. `go vet`/`gofmt` must pass on
every touched file.

### 2. Component abstraction — no change

Nothing to build. Update README's Limitations bullet away instead of
adding code:

> Composition is plain Go functions returning `Node`, plus
> `Fragment(children...) []Node` for grouping pieces — no additional
> component type.

This moves from "limitation" framing to a factual one-liner describing the
model, likely relocated out of Limitations entirely (see below).

### 3. `ScopedStyle` helper

New file `scoped_style.go`:

```go
package html

import "strings"

// ScopedStyle renders a <style> tag whose top-level selectors are each
// prefixed with #scopeID, so the given CSS can't collide with another
// component's rules sharing the page. It is a plain string rewrite, not a
// CSS parser: at-rules (@media, @keyframes, ...) and nested selectors
// inside them are emitted as-is, unprefixed.
func ScopedStyle(scopeID string, css string) Node {
    return rawNode("<style>" + scopeSelectors(scopeID, css) + "</style>")
}

func scopeSelectors(scopeID, css string) string {
    var b strings.Builder
    for {
        open := strings.IndexByte(css, '{')
        if open < 0 {
            b.WriteString(css)
            break
        }
        selectors := css[:open]
        trimmed := strings.TrimSpace(selectors)
        if trimmed == "" || strings.HasPrefix(trimmed, "@") {
            b.WriteString(selectors)
        } else {
            parts := strings.Split(selectors, ",")
            for i, p := range parts {
                if i > 0 {
                    b.WriteString(",")
                }
                b.WriteString(" #" + scopeID + " " + strings.TrimSpace(p))
            }
        }
        b.WriteByte('{')
        close := strings.IndexByte(css, '}')
        if close < 0 {
            b.WriteString(css[open+1:])
            break
        }
        b.WriteString(css[open+1 : close+1])
        css = css[close+1:]
    }
    return b.String()
}
```

Emits raw HTML via `rawNode` (same primitive `UnsafeHTML` uses) since the
CSS body itself is not escaped — same trust model as `<script>` content
today. No new `Node` kind, no change to `render.go`/`validate.go`.

This does not touch `<script>` — there's no non-parser way to scope JS
identifiers. README's bullet gets reworded rather than deleted:

> No JS scoping — `<script>` content is emitted as-is; keeping identifiers
> collision-free across a page is your responsibility. `ScopedStyle`
> namespaces CSS selectors for you; it only rewrites top-level selectors
> (not `@media`/nested rules), so deeply nested or at-rule-wrapped CSS
> still needs manual care.

### Testing

- One `scoped_style_test.go` covering: simple selector, comma-separated
  selectors, an `@media` block passed through unprefixed, and an
  already-scoped output round-tripped through `Render`.
- No new tests per element file — existing `TestRenderTree`-style coverage
  already exercises `globalAttrs`/`Attribute` wiring; new fields follow the
  identical code path `input.go` already proves works.
- `go build ./... && go vet ./... && go test ./...` must stay green
  throughout — this is ~130 files touched, so run it after each batch, not
  only at the end.

### Error handling

None needed beyond what exists: nil-pointer fields are already the
opt-in/opt-out mechanism (`if p.X != nil`), `ScopedStyle` has no failure
mode (it's a pure string transform over caller-supplied text, same trust
boundary as `UnsafeHTML`).

## Out of scope

- No CSS parser, no JS scoping, no schema/generator tool.
- No changes to global attributes (`lang`, `data-*`, `aria-*`, ...) — not
  named in the original Limitations bullets, so adding them now would be
  scope creep beyond what's being closed here.
- No new abstraction beyond `Fragment` for composition.

## README changes

Once implemented, Limitations shrinks to:

```md
## Limitations

- No JS scoping — `<script>` content is emitted as-is; keeping identifiers
  collision-free across a page is your responsibility. `ScopedStyle`
  namespaces CSS selectors for you (see below); it only rewrites top-level
  selectors, so deeply nested or at-rule-wrapped CSS still needs manual
  care.
```

And a new `## Scoped styles` section documents `ScopedStyle` next to
`## Inline scripts`, with a runnable example (verified against real
`go run` output, per this repo's existing README convention).
