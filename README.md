# html

[![Go Reference](https://pkg.go.dev/badge/github.com/nsatyasrikar/html.svg)](https://pkg.go.dev/github.com/nsatyasrikar/html)
[![License: MIT](https://img.shields.io/badge/License-MIT-blue.svg)](LICENSE)

Build HTML trees in Go with typed element constructors, compose them with
ordinary functions, and render them to escaped HTML strings.

Requires Go 1.21+. No third-party dependencies, template files, or code
generation. The compiler checks constructor arguments and field types;
validation checks a small set of tree issues, not HTML conformance.

## Contents

- [Install](#install)
- [Quick start](#quick-start)
- [Elements and props](#elements-and-props)
- [Text and custom attributes](#text-and-custom-attributes)
- [Components, loops, and conditions](#components-loops-and-conditions)
- [Forms](#forms)
- [Images, media, and tables](#images-media-and-tables)
- [Styles](#styles)
- [Scripts and event handlers](#scripts-and-event-handlers)
- [HTTP responses](#http-responses)
- [Validation and warnings](#validation-and-warnings)
- [API reference](#api-reference)
- [Development](#development)
- [License](#license)

## Install

```sh
go get github.com/nsatyasrikar/html
```

Examples describe this checkout. Check your installed version for the newer
attribute fields and scoping helpers.

## Quick start

Save this as `main.go` in a Go module and run `go run .`:

```go
package main

import (
	"fmt"

	html "github.com/nsatyasrikar/html"
)

func main() {
	page := html.Html(html.HtmlProps{
		Attributes: []html.Attribute{{Name: "lang", Value: "en"}},
	},
		html.Head(html.HeadProps{},
			html.Meta(html.MetaProps{
				Attributes: []html.Attribute{{Name: "charset", Value: "utf-8"}},
			}),
			html.Title(html.TitleProps{}, html.Text("Hello")),
		),
		html.Body(html.BodyProps{},
			html.Main(html.MainProps{ID: html.Prop("main")},
				html.H1(html.H1Props{}, html.Text("Hello <world>")),
				html.A(html.AProps{Href: html.Prop("/docs")}, html.Text("Read the docs")),
			),
		),
	)
	fmt.Println("<!doctype html>" + html.Render(page))
}
```

Output:

```html
<!doctype html><html lang="en"><head><meta charset="utf-8"><title>Hello</title></head><body><main id="main"><h1>Hello &lt;world&gt;</h1><a href="/docs">Read the docs</a></main></body></html>
```

`Render` returns compact HTML without indentation or a doctype. Add the
fixed prefix when rendering a complete document. Unless a snippet declares
a function or package, the following Go examples belong inside a function
with the `html` import above. Output examples also use `fmt`; logging uses
`log`.

## Elements and props

Container elements accept children in display order. Void elements such as
`Input`, `Img`, `Br`, `Meta`, and `Link` accept only props. Every constructor
requires its props argument, even when empty.

```go
node := html.Div(html.DivProps{
	ID: html.Prop("intro"), Class: html.Prop("card"),
	Style: html.Prop("padding: 1rem"), Title: html.Prop("Introduction"),
}, html.Text("Welcome"), html.Br(html.BrProps{}), html.Text("Start here"))
_ = html.Render(node)
```

Common fields are `ID`, `Class`, `Style`, `Title` (`*string`) and
`Attributes` (`[]html.Attribute`). Element-specific fields include `Href`
on `A`, `Src` on `Img`, and `Rows` on `Textarea`. Coverage varies; use
`Attributes` for fields that are not exposed. Numeric HTML values such as
dimensions and row counts are strings.

`Prop[T](value T) *T` takes a pointer to a value. A nil string prop omits
the attribute; `Prop("")` renders an empty value.

### Boolean presence

Named boolean props currently use pointer presence: both `Prop(true)` and
`Prop(false)` emit the attribute. Use `nil` to omit it. For a dynamic flag:

```go
disabled := false
var disabledProp *bool
if disabled {
	disabledProp = html.Prop(true)
}
node := html.Button(html.ButtonProps{Disabled: disabledProp}, html.Text("Save"))
fmt.Println(html.Render(node))
// <button>Save</button>
```

Custom boolean attributes instead inspect `Value`: only the exact string
`"true"` emits a bare attribute when `Boolean` is true.

## Text and custom attributes

`Text` escapes text content. Attribute values are escaped during rendering.
Custom names must match `^[A-Za-z_:][A-Za-z0-9_.:-]*$`; invalid names are
silently discarded.

```go
node := html.Div(html.DivProps{
	Attributes: []html.Attribute{
		{Name: "data-testid", Value: "hero"},
		{Name: "aria-label", Value: "Tips & tricks"},
		{Name: "hidden", Value: "true", Boolean: true},
	},
}, html.Text("Use <Enter> & continue"))
fmt.Println(html.Render(node))
// <div data-testid="hero" aria-label="Tips &amp; tricks" hidden>Use &lt;Enter&gt; &amp; continue</div>
```

Custom attributes follow common attributes and precede element-specific
ones. They do not override named fields: setting `ID` and a custom `id`
renders both and produces a duplicate-attribute warning.

For trusted, already-rendered markup:

```go
node := html.Div(html.DivProps{}, html.UnsafeHTML("<strong>Trusted markup</strong>"))
fmt.Println(html.Render(node))
// <div><strong>Trusted markup</strong></div>
```

`UnsafeHTML` emits input verbatim and produces an `unsafe-html` warning.
It does not sanitize HTML. Escaping attribute values does not validate URL
schemes or make JavaScript safe; choose those values separately.

## Components, loops, and conditions

Components are ordinary functions returning `Node`. `Fragment` returns
`[]Node` for grouping siblings without adding a wrapper. Expand the slice
with `...` when passing it to a constructor.

```go
func Card(title string, items []string, showHelp bool) html.Node {
	rows := make([]html.Node, 0, len(items))
	for _, item := range items {
		rows = append(rows, html.Li(html.LiProps{}, html.Text(item)))
	}
	children := html.Fragment(
		html.H2(html.H2Props{}, html.Text(title)),
		html.Ul(html.UlProps{}, rows...),
	)
	if showHelp {
		children = append(children, html.P(html.PProps{}, html.Text("Choose an item.")))
	}
	return html.Article(html.ArticleProps{Class: html.Prop("card")}, children...)
}
```

Call `html.Render(Card("Tasks", []string{"Read", "Build"}, true))`.
An empty children slice is valid; omit optional nodes instead of inserting
`nil`. Nil roots and nil children are not supported.

## Forms

```go
form := html.Form(html.FormProps{Action: html.Prop("/signup"), Method: html.Prop("post")},
	html.Label(html.LabelProps{For: html.Prop("email")}, html.Text("Email")),
	html.Input(html.InputProps{
		ID: html.Prop("email"), Name: html.Prop("email"),
		Type: html.Prop(html.InputTypeEmail), Required: html.Prop(true),
		Placeholder: html.Prop("you@example.com"),
	}),
	html.Label(html.LabelProps{For: html.Prop("plan")}, html.Text("Plan")),
	html.Select(html.SelectProps{ID: html.Prop("plan"), Name: html.Prop("plan")},
		html.Optgroup(html.OptgroupProps{Label: html.Prop("Subscriptions")},
			html.Option(html.OptionProps{Value: html.Prop("starter"), Selected: html.Prop(true)}, html.Text("Starter")),
			html.Option(html.OptionProps{Value: html.Prop("team")}, html.Text("Team")),
		),
	),
	html.Label(html.LabelProps{For: html.Prop("bio")}, html.Text("About you")),
	html.Textarea(html.TextareaProps{
		ID: html.Prop("bio"), Name: html.Prop("bio"), Rows: html.Prop("4"), Cols: html.Prop("40"),
	}, html.Text("I build things.")),
	html.Button(html.ButtonProps{
		Attributes: []html.Attribute{{Name: "type", Value: "submit"}},
	}, html.Text("Sign up")),
)
_ = html.Render(form)
```

`InputTypeText`, `InputTypeEmail`, `InputTypePassword`, and `InputTypeNumber`
are predefined. Other values can be expressed as `html.InputType("date")`.
The type is a string type, not an enum validator.

`FieldsetProps.Disabled` is currently not wired into its constructor; use a
custom `disabled` boolean attribute on `Fieldset` when needed.

## Images, media, and tables

```go
image := html.Img(html.ImgProps{
	Src: html.Prop("/photo.jpg"), Alt: html.Prop("A mountain lake"),
	Width: html.Prop("640"), Height: html.Prop("480"), Loading: html.Prop("lazy"),
	Srcset: html.Prop("/photo.jpg 1x, /photo@2x.jpg 2x"),
})
video := html.Video(html.VideoProps{Controls: html.Prop(true), Poster: html.Prop("/poster.jpg")},
	html.Source(html.SourceProps{Src: html.Prop("/tour.mp4"), Type: html.Prop("video/mp4")}),
	html.Track(html.TrackProps{
		Src: html.Prop("/captions.vtt"), Kind: html.Prop("captions"),
		Srclang: html.Prop("en"), Label: html.Prop("English"), Default: html.Prop(true),
	}),
	html.Text("Your browser does not support video."),
)
table := html.Table(html.TableProps{},
	html.Caption(html.CaptionProps{}, html.Text("Subscriptions")),
	html.Thead(html.TheadProps{}, html.Tr(html.TrProps{},
		html.Th(html.ThProps{Scope: html.Prop("col")}, html.Text("Plan")),
		html.Th(html.ThProps{Scope: html.Prop("col")}, html.Text("Seats")),
	)),
	html.Tbody(html.TbodyProps{}, html.Tr(html.TrProps{},
		html.Td(html.TdProps{}, html.Text("Team")),
		html.Td(html.TdProps{}, html.Text("5")),
	)),
)
_ = html.Render(html.Section(html.SectionProps{}, image, video, table))
```

Other constructors follow the same pattern: `Audio`, `Picture`, `Canvas`,
`Iframe`, `Details`/`Summary`, `Dialog`, `Meter`, `Progress`, and semantic
elements such as `Time`, `Blockquote`, and `Data`. Consult each element's
props definition for supported fields.

## Styles

For external CSS, use `Link`. For inline CSS, `Style` supplies the tag and
`UnsafeHTML` supplies trusted CSS without HTML text escaping:

```go
external := html.Link(html.LinkProps{Rel: html.Prop("stylesheet"), Href: html.Prop("/app.css")})
inline := html.Style(html.StyleProps{Media: html.Prop("print")},
	html.UnsafeHTML(".toolbar { display: none; }"),
)
_ = html.Render(html.Head(html.HeadProps{}, external, inline))
```

`ScopedStyle(scopeID, css)` supplies its own style tag and prefixes simple
selectors with `#scopeID`. Put targets inside an element with that id;
the prefix selects descendants, not the wrapper itself.

```go
card := html.Div(html.DivProps{ID: html.Prop("card-1")},
	html.ScopedStyle("card-1", ".title { color: red; }"),
	html.H2(html.H2Props{Class: html.Prop("title")}, html.Text("Hello")),
)
fmt.Println(html.Render(card))
// <div id="card-1"><style> #card-1 .title{ color: red; }</style><h2 class="title">Hello</h2></div>
```

This is a string rewrite, not a CSS parser. At-rules and their contents
remain unscoped. Commas in functional selectors, braces in strings or
comments, and complex nesting are not reliably handled. Use simple trusted
rules and an id that is already a valid CSS identifier. For complex CSS,
author explicit scoped selectors in a stylesheet. `ScopedStyle` does not
escape closing style tags and emits an `unsafe-html` warning.

## Scripts and event handlers

Use `ScriptTag` for a script element. External scripts support `Src`,
`Type`, `Async`, and `Defer`; inline scripts use `ScriptCode` for the body:

```go
external := html.ScriptTag(html.ScriptTagProps{Src: html.Prop("/app.js"), Defer: html.Prop(true)})
inline := html.ScriptTag(html.ScriptTagProps{},
	html.ScriptCode("document.title = 'Ready';"),
)
button := html.Button(html.ButtonProps{
	OnClick: html.Prop(html.Script("console.log('clicked')")),
}, html.Text("Try it"))
_ = html.Render(html.Div(html.DivProps{}, external, inline, button))
```

`Script` is the string type used by `ButtonProps.OnClick` and
`InputProps.OnInput`/`OnChange`. Event attributes are HTML-escaped; the
JavaScript itself must be trusted.

`ScriptCode` emits raw JavaScript and rewrites case-insensitive `</script`
sequences as `<\/script` so they cannot terminate the surrounding tag. It
does not add a tag, sanitize JavaScript, or emit an `unsafe-html` warning.
Use it inside `ScriptTag`, not `Text`, which would HTML-escape the code.

`ScopedScript` adds a script tag and an immediately invoked function:

```go
fmt.Println(html.Render(html.ScopedScript("var count = 0;")))
// <script>(function(){
// var count = 0;
// })();</script>
```

The wrapper keeps ordinary declarations local. It is not a sandbox:
explicit global writes and DOM changes remain possible. It accepts a
function body, not module-level `import`/`export`. Unlike `ScriptCode`,
it does not escape script end tags; use trusted source without literal
closing script tags. It emits an `unsafe-html` warning.

## HTTP responses

`Render` returns a string, so it works with `net/http` without an adapter:

```go
package main

import (
	"io"
	"log"
	"net/http"

	html "github.com/nsatyasrikar/html"
)

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		page := html.Html(html.HtmlProps{},
			html.Head(html.HeadProps{}, html.Title(html.TitleProps{}, html.Text("Welcome"))),
			html.Body(html.BodyProps{}, html.H1(html.H1Props{}, html.Text("Hello, "+r.URL.Query().Get("name")))),
		)
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		if _, err := io.WriteString(w, "<!doctype html>"+html.Render(page)); err != nil {
			log.Print(err)
		}
	})
	log.Fatal(http.ListenAndServe("localhost:8080", nil))
}
```

Run this program and visit `http://localhost:8080/?name=Go`.
Rendering buffers the whole tree into a string; there is no streaming API.

## Validation and warnings

`Validate(node)` returns `[]Warning` without invoking the warning handler.
`Render(node)` always validates, reports warnings, and continues rendering.
Warnings do not repair or remove offending markup.

```go
node := html.Div(html.DivProps{ID: html.Prop("same")},
	html.Span(html.SpanProps{ID: html.Prop("same")}, html.Text("Repeated id")),
)
for _, w := range html.Validate(node) {
	fmt.Printf("%s: tag=%s path=%s\n", w.Code, w.Tag, w.Path)
}
// duplicate-id: tag=span path=/div[0]
```

| Code | Meaning |
| --- | --- |
| `duplicate-attribute` | An attribute name occurs more than once on an element. |
| `duplicate-id` | A nonempty id is reused in the tree. |
| `void-children` | An internally marked void element contains children; public void constructors prevent this. |
| `unsafe-html` | A raw HTML node is present, including either scoping helper. |
| `deprecated` | An internal element node is marked deprecated; current constructors do not set this flag. |

`Warning` has string fields `Code`, `Message`, `Tag`, and `Path`. Paths use
zero-based child indexes; the root path is empty. Raw nodes are opaque:
their markup is not parsed. Validation does not check nesting rules,
required attributes, URLs, CSS, or JavaScript.

The default handler logs warnings. Set a process-wide handler at startup:

```go
html.SetWarningHandler(func(w html.Warning) {
	log.Printf("[%s] %s at %s", w.Code, w.Message, w.Path)
})
```

`html.SetWarningHandler(nil)` disables reporting, but still runs validation.
Handlers may be called concurrently by concurrent renders. For strict
behavior, inspect `Validate` and decide which codes should prevent your
application from calling `Render`.

## API reference

| API | Purpose |
| --- | --- |
| `Node` | Sealed interface returned by constructors; compose through functions. |
| `Div(DivProps, ...Node)` and other elements | Build trees with typed props. |
| `Prop[T](T) *T` | Supply optional property values. |
| `Attribute{Name, Value, Boolean}` | Supply custom attributes. |
| `Text(string) Node` | Emit escaped text. |
| `UnsafeHTML(string) Node` | Emit trusted markup verbatim. |
| `Fragment(...Node) []Node` | Group siblings for slice expansion. |
| `Script` / `ScriptCode(string) Node` | Supply event code / inline script bodies. |
| `ScopedStyle(string, string) Node` | Prefix simple CSS selectors and emit a style tag. |
| `ScopedScript(string) Node` | Wrap trusted code in an IIFE and script tag. |
| `Render(Node) string` | Validate, report warnings, and render HTML. |
| `Validate(Node) []Warning` | Inspect warnings without rendering. |
| `WarningHandler` / `SetWarningHandler(WarningHandler)` | Configure reporting. |

Browse the [package reference](https://pkg.go.dev/github.com/nsatyasrikar/html)
for published signatures. For your checkout, run `go doc .` or
`go doc . VideoProps`. Element implementations live in individual files,
such as [input.go](input.go), [video.go](video.go), and [table.go](table.go).
The script element is named `ScriptTag` because `Script` is the event-code
type.

## Development

```sh
go build ./...
go vet ./...
go test ./...
gofmt -l .
```

No network access or API key is needed to build or test once Go is installed.
See [Architecture](notes/Architecture.md) and the
[Decisions Log](notes/Decisions%20Log.md) for implementation context.

## License

[MIT](LICENSE)
