# html

[![Go Reference](https://pkg.go.dev/badge/github.com/nsatyasrikar/html.svg)](https://pkg.go.dev/github.com/nsatyasrikar/html)
[![Go Report Card](https://goreportcard.com/badge/github.com/nsatyasrikar/html)](https://goreportcard.com/report/github.com/nsatyasrikar/html)
[![License: MIT](https://img.shields.io/badge/License-MIT-blue.svg)](LICENSE)

A typed Go HTML builder. Every element (`Div`, `P`, `Input`, ...) is a plain
function taking a `*Props` struct and returning a `Node`; `Render` walks the
tree to an escaped HTML string.

No template files, no string concatenation, no runtime parsing — the tree
you build is the tree you render, checked by the Go compiler along the way.

## Install

```sh
go get github.com/nsatyasrikar/html
```

Requires Go 1.21+. Zero third-party dependencies.

## Contents

- [Quick start](#quick-start)
- [Props are pointers](#props-are-pointers)
- [Text vs raw HTML](#text-vs-raw-html)
- [Custom and boolean attributes](#custom-and-boolean-attributes)
- [Inline scripts](#inline-scripts)
- [Validation warnings](#validation-warnings)
- [Code generation](#code-generation)
- [Building and testing](#building-and-testing)
- [Limitations](#limitations)
- [License](#license)

## Quick start

```go
package main

import (
	"fmt"

	"github.com/nsatyasrikar/html"
)

func main() {
	page := html.Html(html.HtmlProps{},
		html.Head(html.HeadProps{},
			html.Title(html.TitleProps{}, html.Text("Hello")),
		),
		html.Body(html.BodyProps{},
			html.H1(html.H1Props{Class: html.Prop("heading")},
				html.Text("Hello <world>"),
			),
		),
	)
	fmt.Println(html.Render(page))
}
```

```html
<html><head><title>Hello</title></head><body><h1 class="heading">Hello &lt;world&gt;</h1></body></html>
```

Text content is HTML-escaped automatically (`<world>` → `&lt;world&gt;`).

## Props are pointers

Every prop is a pointer so "unset" and "set to the zero value" stay
distinguishable — an unset `*bool` renders nothing, `Prop(false)` still
renders `false`-driven markup where it matters (e.g. boolean attributes only
appear when the value is `true`). `Prop` is a generic helper for the common
case of taking a pointer to a literal:

```go
func Prop[T any](value T) *T
```

```go
html.Input(html.InputProps{
	Type:     html.Prop(html.InputTypeEmail),
	Disabled: html.Prop(true),
})
// <input type="email" disabled>
```

Void elements (`input`, `img`, `br`, `hr`, `meta`, `link`, ...) take no
`children ...Node` parameter — that's a compile-time guarantee, not a
runtime check.

## Text vs raw HTML

```go
html.Text("Hello <world>")   // escaped: Hello &lt;world&gt;
html.UnsafeHTML("<b>bold</b>") // emitted verbatim, no escaping
```

`UnsafeHTML` is for markup you already trust (e.g. output from a Markdown
renderer). Every `UnsafeHTML` node raises an `"unsafe-html"` warning from
`Validate` so it's easy to grep for the places raw HTML enters the tree.

## Custom and boolean attributes

`Attributes []Attribute` on every `Props` struct covers anything not
already a named field:

```go
html.Div(html.DivProps{
	Attributes: []Attribute{
		{Name: "data-testid", Value: "hero"},
		{Name: "hidden", Boolean: true, Value: "true"},
	},
})
// <div data-testid="hero" hidden>
```

Attribute names are checked against `^[A-Za-z_:][A-Za-z0-9_.:-]*$`
(`attribute.go`) before rendering — anything that doesn't match a valid
HTML attribute name is silently dropped, not rendered.

## Inline scripts

```go
html.ScriptTag(html.ScriptTagProps{},
	html.ScriptCode(`document.title = "hi"`),
)
```

`ScriptCode` escapes any literal `</script` sequence inside the source
(`<\/script`) so untrusted-looking script bodies can't prematurely close the
`<script>` tag.

## Validation warnings

`Render` always runs `Validate` first and reports every warning through the
package-level handler before writing output:

```go
html.SetWarningHandler(func(w html.Warning) {
	log.Printf("[%s] %s at %s", w.Code, w.Message, w.Path)
})
```

Default handler logs via `log.Printf`. Checks run in one tree walk: void
elements with children (`void-children`), duplicate attributes on the same
tag (`duplicate-attribute`), duplicate `id` values anywhere in the tree
(`duplicate-id`), and any `UnsafeHTML` usage (`unsafe-html`). Warnings never
block rendering — call `Validate` directly if you need to fail on them.

## Building and testing

```sh
go build ./...
go test ./...
```

No API key or network access required to build or test.

## Scoped styles and scripts

`ScopedStyle` prefixes selectors in ordinary CSS rules with an id:

```go
html.Render(html.ScopedStyle("card-1", ".title { color: red; }"))
// <style> #card-1 .title { color: red; }</style>
```

`ScopedScript` wraps a script body in an IIFE so its declarations stay local:

```go
html.Render(html.ScopedScript("var count = 0;"))
// <script>(function(){
// var count = 0;
// })();</script>
```

## License

[MIT](LICENSE)
