# Close README Limitations Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Hand-write real HTML attributes onto every element file that has any, add `ScopedStyle`/`ScopedScript` helpers to close CSS/JS collision risk, and delete README's `## Limitations` section entirely.

**Architecture:** Every element constructor already follows one template (`globalAttrs` + conditional `append(a, Attribute{...})` + `elementNode{...}`) established by `input.go`/`a.go`/`meta.go`/`link.go`. This plan repeats that template on ~50 more files, grouped into tasks by element category so each task's diff is reviewable as one unit. Two new standalone files (`scoped_style.go`, `scoped_script.go`) add the scoping helpers. No schema, no generator, no new `Node` kind.

**Tech Stack:** Go 1.21+, stdlib only (`strings` for `scoped_style.go`), package `html` at repo root.

**Spec:** `docs/superpowers/specs/2026-09-09-close-limitations-design.md`

## Global Constraints

- No `cmd/` tool, no schema file (`index.json` or otherwise) gets reintroduced — every file stays independently hand-editable, per `notes/Decisions Log.md#htmlgen removal`.
- Boolean attributes use the existing convention: `Attribute{name, "true", true}` (third field marks it boolean; `render.go` emits the bare name with no value).
- Every touched file must pass `gofmt -l` (no diff) and the whole module must pass `go build ./... && go vet ./... && go test ./...` after each task.
- Elements with no real attributes beyond the generic four (`ID`/`Class`/`Style`/`Title`) are left untouched — do not add unused fields.
- Commit after each task, not once at the end.

---

## Task 1: `ScopedStyle` helper

**Files:**
- Create: `scoped_style.go`
- Test: `scoped_style_test.go`

**Interfaces:**
- Produces: `func ScopedStyle(scopeID string, css string) Node` — later tasks/README examples call this directly.

- [ ] **Step 1: Write the failing test**

```go
package html

import (
	"strings"
	"testing"
)

func TestScopedStyleSimpleSelector(t *testing.T) {
	out := Render(ScopedStyle("card-1", ".title { color: red; }"))
	if !strings.Contains(out, "#card-1 .title") {
		t.Fatalf("expected scoped selector in output, got: %s", out)
	}
}

func TestScopedStyleCommaSeparatedSelectors(t *testing.T) {
	out := Render(ScopedStyle("card-1", ".a, .b { color: red; }"))
	if !strings.Contains(out, "#card-1 .a") || !strings.Contains(out, "#card-1 .b") {
		t.Fatalf("expected both selectors scoped, got: %s", out)
	}
}

func TestScopedStyleAtRulePassthrough(t *testing.T) {
	css := "@media (min-width: 600px) { .a { color: red; } }"
	out := Render(ScopedStyle("card-1", css))
	if !strings.Contains(out, "@media (min-width: 600px)") {
		t.Fatalf("expected @media passed through unprefixed, got: %s", out)
	}
}

func TestScopedStyleRendersStyleTag(t *testing.T) {
	out := Render(ScopedStyle("card-1", ".a { color: red; }"))
	if !strings.HasPrefix(out, "<style>") || !strings.HasSuffix(out, "</style>") {
		t.Fatalf("expected wrapped in <style> tags, got: %s", out)
	}
}
```

- [ ] **Step 2: Run test to verify it fails**

Run: `go test -run TestScopedStyle ./...`
Expected: FAIL with "undefined: ScopedStyle"

- [ ] **Step 3: Write minimal implementation**

```go
package html

import "strings"

// ScopedStyle renders a <style> tag whose top-level selectors are each
// prefixed with #scopeID, so the given CSS can't collide with another
// component's rules sharing the page. It is a plain string rewrite, not a
// CSS parser: at-rules (@media, @keyframes, ...) and their nested
// selectors are emitted as-is, unprefixed.
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
		css = css[open+1:]
		close := strings.IndexByte(css, '}')
		if close < 0 {
			b.WriteString(css)
			break
		}
		b.WriteString(css[:close+1])
		css = css[close+1:]
	}
	return b.String()
}
```

- [ ] **Step 4: Run test to verify it passes**

Run: `go test -run TestScopedStyle ./...`
Expected: PASS (all 4 subtests)

- [ ] **Step 5: Commit**

```bash
git add scoped_style.go scoped_style_test.go
git commit -m "add ScopedStyle helper for collision-free component CSS"
```

---

## Task 2: `ScopedScript` helper

**Files:**
- Create: `scoped_script.go`
- Test: `scoped_script_test.go`

**Interfaces:**
- Produces: `func ScopedScript(js string) Node` — later README example calls this directly.

- [ ] **Step 1: Write the failing test**

```go
package html

import (
	"strings"
	"testing"
)

func TestScopedScriptWrapsInIIFE(t *testing.T) {
	out := Render(ScopedScript("var x = 1;"))
	if !strings.Contains(out, "(function(){") || !strings.Contains(out, "})();") {
		t.Fatalf("expected IIFE wrapper, got: %s", out)
	}
}

func TestScopedScriptBodyNotBareGlobal(t *testing.T) {
	out := Render(ScopedScript("var leaked = 1;"))
	// "var leaked" must appear only inside the wrapper, never as
	// "</script>var leaked" or similar outside it.
	closeIdx := strings.Index(out, "})();")
	if closeIdx < 0 {
		t.Fatalf("wrapper close not found: %s", out)
	}
	if strings.Contains(out[closeIdx:], "var leaked") {
		t.Fatalf("declaration leaked outside IIFE: %s", out)
	}
}

func TestScopedScriptRendersScriptTag(t *testing.T) {
	out := Render(ScopedScript("var x = 1;"))
	if !strings.HasPrefix(out, "<script>") || !strings.HasSuffix(out, "</script>") {
		t.Fatalf("expected wrapped in <script> tags, got: %s", out)
	}
}
```

- [ ] **Step 2: Run test to verify it fails**

Run: `go test -run TestScopedScript ./...`
Expected: FAIL with "undefined: ScopedScript"

- [ ] **Step 3: Write minimal implementation**

```go
package html

// ScopedScript renders a <script> tag whose body runs inside an IIFE, so
// every var/let/const/function it declares is invisible outside — no
// parser needed, JavaScript's own function scoping does the isolating.
func ScopedScript(js string) Node {
	return rawNode("<script>(function(){\n" + js + "\n})();</script>")
}
```

- [ ] **Step 4: Run test to verify it passes**

Run: `go test -run TestScopedScript ./...`
Expected: PASS (all 3 subtests)

- [ ] **Step 5: Commit**

```bash
git add scoped_script.go scoped_script_test.go
git commit -m "add ScopedScript helper for collision-free component JS"
```

---

## Task 3: Void embed/media-source elements (area, base, col, embed, img, param, source, track, hr)

**Files:**
- Modify: `area.go`, `base.go`, `col.go`, `embed.go`, `img.go`, `param.go`, `source.go`, `track.go`, `hr.go`
- Test: `element_attrs_test.go` (new shared test file for Tasks 3-10; each task appends to it)

**Interfaces:**
- Produces: `AreaProps{Href, Alt, Coords, Shape, Target}`, `BaseProps{Href, Target}`, `ColProps{Span}`, `EmbedProps{Src, Type, Width, Height}`, `ImgProps{..., Width, Height, Srcset, Sizes, Loading}` (existing `Src`/`Alt` untouched), `ParamProps{Name, Value}`, `SourceProps{Src, Type, Srcset, Sizes, Media}`, `TrackProps{Src, Kind, Srclang, Label, Default *bool}`, `HrProps{Color, Size, Width}`.

- [ ] **Step 1: Write the failing test**

```go
package html

import (
	"strings"
	"testing"
)

func TestAreaAttrs(t *testing.T) {
	href, alt, coords, shape, target := "#x", "map area", "0,0,10,10", "rect", "_blank"
	out := Render(Area(AreaProps{Href: &href, Alt: &alt, Coords: &coords, Shape: &shape, Target: &target}))
	for _, want := range []string{`href="#x"`, `alt="map area"`, `coords="0,0,10,10"`, `shape="rect"`, `target="_blank"`} {
		if !strings.Contains(out, want) {
			t.Fatalf("expected %s in %s", want, out)
		}
	}
}

func TestBaseAttrs(t *testing.T) {
	href, target := "/", "_self"
	out := Render(Base(BaseProps{Href: &href, Target: &target}))
	if !strings.Contains(out, `href="/"`) || !strings.Contains(out, `target="_self"`) {
		t.Fatalf("got: %s", out)
	}
}

func TestColAttrs(t *testing.T) {
	span := "2"
	out := Render(Col(ColProps{Span: &span}))
	if !strings.Contains(out, `span="2"`) {
		t.Fatalf("got: %s", out)
	}
}

func TestEmbedAttrs(t *testing.T) {
	src, typ, w, h := "movie.swf", "application/x-shockwave-flash", "400", "300"
	out := Render(Embed(EmbedProps{Src: &src, Type: &typ, Width: &w, Height: &h}))
	for _, want := range []string{`src="movie.swf"`, `type="application/x-shockwave-flash"`, `width="400"`, `height="300"`} {
		if !strings.Contains(out, want) {
			t.Fatalf("expected %s in %s", want, out)
		}
	}
}

func TestImgNewAttrs(t *testing.T) {
	w, h, srcset, sizes, loading := "600", "400", "a.jpg 1x, b.jpg 2x", "(min-width: 600px) 600px", "lazy"
	out := Render(Img(ImgProps{Width: &w, Height: &h, Srcset: &srcset, Sizes: &sizes, Loading: &loading}))
	for _, want := range []string{`width="600"`, `height="400"`, `srcset="a.jpg 1x, b.jpg 2x"`, `sizes="(min-width: 600px) 600px"`, `loading="lazy"`} {
		if !strings.Contains(out, want) {
			t.Fatalf("expected %s in %s", want, out)
		}
	}
}

func TestParamAttrs(t *testing.T) {
	name, value := "autoplay", "true"
	out := Render(Param(ParamProps{Name: &name, Value: &value}))
	if !strings.Contains(out, `name="autoplay"`) || !strings.Contains(out, `value="true"`) {
		t.Fatalf("got: %s", out)
	}
}

func TestSourceAttrs(t *testing.T) {
	src, typ, srcset, sizes, media := "vid.mp4", "video/mp4", "vid-2x.mp4 2x", "(min-width: 600px)", "(min-width: 600px)"
	out := Render(Source(SourceProps{Src: &src, Type: &typ, Srcset: &srcset, Sizes: &sizes, Media: &media}))
	for _, want := range []string{`src="vid.mp4"`, `type="video/mp4"`, `srcset="vid-2x.mp4 2x"`} {
		if !strings.Contains(out, want) {
			t.Fatalf("expected %s in %s", want, out)
		}
	}
}

func TestTrackAttrs(t *testing.T) {
	src, kind, srclang, label := "captions.vtt", "captions", "en", "English"
	def := true
	out := Render(Track(TrackProps{Src: &src, Kind: &kind, Srclang: &srclang, Label: &label, Default: &def}))
	for _, want := range []string{`src="captions.vtt"`, `kind="captions"`, `srclang="en"`, `label="English"`, "default"} {
		if !strings.Contains(out, want) {
			t.Fatalf("expected %s in %s", want, out)
		}
	}
}

func TestHrAttrs(t *testing.T) {
	color, size, width := "red", "2", "80%"
	out := Render(Hr(HrProps{Color: &color, Size: &size, Width: &width}))
	for _, want := range []string{`color="red"`, `size="2"`, `width="80%"`} {
		if !strings.Contains(out, want) {
			t.Fatalf("expected %s in %s", want, out)
		}
	}
}
```

- [ ] **Step 2: Run test to verify it fails**

Run: `go test -run 'TestArea|TestBase|TestCol|TestEmbed|TestImgNewAttrs|TestParam|TestSource|TestTrack|TestHr' ./...`
Expected: FAIL (unknown fields `Href`, `Span`, etc.)

- [ ] **Step 3: Write minimal implementation**

`area.go`:
```go
package html

type AreaProps struct {
	ID, Class, Style, Title, Href, Alt, Coords, Shape, Target *string
	Attributes                                                 []Attribute
}

func Area(p AreaProps) Node {
	a := globalAttrs(p.ID, p.Class, p.Style, p.Title, p.Attributes)
	if p.Href != nil {
		a = append(a, Attribute{"href", *p.Href, false})
	}
	if p.Alt != nil {
		a = append(a, Attribute{"alt", *p.Alt, false})
	}
	if p.Coords != nil {
		a = append(a, Attribute{"coords", *p.Coords, false})
	}
	if p.Shape != nil {
		a = append(a, Attribute{"shape", *p.Shape, false})
	}
	if p.Target != nil {
		a = append(a, Attribute{"target", *p.Target, false})
	}
	return elementNode{name: "area", props: a, void: true}
}
```

`base.go`:
```go
package html

type BaseProps struct {
	ID, Class, Style, Title, Href, Target *string
	Attributes                            []Attribute
}

func Base(p BaseProps) Node {
	a := globalAttrs(p.ID, p.Class, p.Style, p.Title, p.Attributes)
	if p.Href != nil {
		a = append(a, Attribute{"href", *p.Href, false})
	}
	if p.Target != nil {
		a = append(a, Attribute{"target", *p.Target, false})
	}
	return elementNode{name: "base", props: a, void: true}
}
```

`col.go`:
```go
package html

type ColProps struct {
	ID, Class, Style, Title, Span *string
	Attributes                    []Attribute
}

func Col(p ColProps) Node {
	a := globalAttrs(p.ID, p.Class, p.Style, p.Title, p.Attributes)
	if p.Span != nil {
		a = append(a, Attribute{"span", *p.Span, false})
	}
	return elementNode{name: "col", props: a, void: true}
}
```

`embed.go`:
```go
package html

type EmbedProps struct {
	ID, Class, Style, Title, Src, Type, Width, Height *string
	Attributes                                        []Attribute
}

func Embed(p EmbedProps) Node {
	a := globalAttrs(p.ID, p.Class, p.Style, p.Title, p.Attributes)
	if p.Src != nil {
		a = append(a, Attribute{"src", *p.Src, false})
	}
	if p.Type != nil {
		a = append(a, Attribute{"type", *p.Type, false})
	}
	if p.Width != nil {
		a = append(a, Attribute{"width", *p.Width, false})
	}
	if p.Height != nil {
		a = append(a, Attribute{"height", *p.Height, false})
	}
	return elementNode{name: "embed", props: a, void: true}
}
```

`img.go`:
```go
package html

type ImgProps struct {
	ID, Class, Style, Title, Src, Alt, Width, Height, Srcset, Sizes, Loading *string
	Attributes                                                               []Attribute
}

func Img(p ImgProps) Node {
	a := globalAttrs(p.ID, p.Class, p.Style, p.Title, p.Attributes)
	if p.Src != nil {
		a = append(a, Attribute{"src", *p.Src, false})
	}
	if p.Alt != nil {
		a = append(a, Attribute{"alt", *p.Alt, false})
	}
	if p.Width != nil {
		a = append(a, Attribute{"width", *p.Width, false})
	}
	if p.Height != nil {
		a = append(a, Attribute{"height", *p.Height, false})
	}
	if p.Srcset != nil {
		a = append(a, Attribute{"srcset", *p.Srcset, false})
	}
	if p.Sizes != nil {
		a = append(a, Attribute{"sizes", *p.Sizes, false})
	}
	if p.Loading != nil {
		a = append(a, Attribute{"loading", *p.Loading, false})
	}
	return elementNode{name: "img", props: a, void: true}
}
```

`param.go`:
```go
package html

type ParamProps struct {
	ID, Class, Style, Title, Name, Value *string
	Attributes                           []Attribute
}

func Param(p ParamProps) Node {
	a := globalAttrs(p.ID, p.Class, p.Style, p.Title, p.Attributes)
	if p.Name != nil {
		a = append(a, Attribute{"name", *p.Name, false})
	}
	if p.Value != nil {
		a = append(a, Attribute{"value", *p.Value, false})
	}
	return elementNode{name: "param", props: a, void: true}
}
```

`source.go`:
```go
package html

type SourceProps struct {
	ID, Class, Style, Title, Src, Type, Srcset, Sizes, Media *string
	Attributes                                                []Attribute
}

func Source(p SourceProps) Node {
	a := globalAttrs(p.ID, p.Class, p.Style, p.Title, p.Attributes)
	if p.Src != nil {
		a = append(a, Attribute{"src", *p.Src, false})
	}
	if p.Type != nil {
		a = append(a, Attribute{"type", *p.Type, false})
	}
	if p.Srcset != nil {
		a = append(a, Attribute{"srcset", *p.Srcset, false})
	}
	if p.Sizes != nil {
		a = append(a, Attribute{"sizes", *p.Sizes, false})
	}
	if p.Media != nil {
		a = append(a, Attribute{"media", *p.Media, false})
	}
	return elementNode{name: "source", props: a, void: true}
}
```

`track.go`:
```go
package html

type TrackProps struct {
	ID, Class, Style, Title, Src, Kind, Srclang, Label *string
	Default                                            *bool
	Attributes                                         []Attribute
}

func Track(p TrackProps) Node {
	a := globalAttrs(p.ID, p.Class, p.Style, p.Title, p.Attributes)
	if p.Src != nil {
		a = append(a, Attribute{"src", *p.Src, false})
	}
	if p.Kind != nil {
		a = append(a, Attribute{"kind", *p.Kind, false})
	}
	if p.Srclang != nil {
		a = append(a, Attribute{"srclang", *p.Srclang, false})
	}
	if p.Label != nil {
		a = append(a, Attribute{"label", *p.Label, false})
	}
	if p.Default != nil {
		a = append(a, Attribute{"default", "true", true})
	}
	return elementNode{name: "track", props: a, void: true}
}
```

`hr.go`:
```go
package html

type HrProps struct {
	ID, Class, Style, Title, Color, Size, Width *string
	Attributes                                  []Attribute
}

func Hr(p HrProps) Node {
	a := globalAttrs(p.ID, p.Class, p.Style, p.Title, p.Attributes)
	if p.Color != nil {
		a = append(a, Attribute{"color", *p.Color, false})
	}
	if p.Size != nil {
		a = append(a, Attribute{"size", *p.Size, false})
	}
	if p.Width != nil {
		a = append(a, Attribute{"width", *p.Width, false})
	}
	return elementNode{name: "hr", props: a, void: true}
}
```

- [ ] **Step 4: Run test to verify it passes**

Run: `go test ./...`
Expected: PASS (all tests, no regressions on existing `Img` `Src`/`Alt` tests)

- [ ] **Step 5: Commit**

```bash
git add area.go base.go col.go embed.go img.go param.go source.go track.go hr.go element_attrs_test.go
git commit -m "extend void media/embed elements with real attrs"
```

---

## Task 4: Inline semantic elements (blockquote, del, ins, q, time, data, li, ol, details, dialog)

**Files:**
- Modify: `blockquote.go`, `del.go`, `ins.go`, `q.go`, `time.go`, `data.go`, `li.go`, `ol.go`, `details.go`, `dialog.go`
- Test: append to `element_attrs_test.go`

**Interfaces:**
- Produces: `BlockquoteProps{Cite}`, `DelProps{Cite, Datetime}`, `InsProps{Cite, Datetime}`, `QProps{Cite}`, `TimeProps{Datetime}`, `DataProps{Value}`, `LiProps{Value}`, `OlProps{Start, Type, Reversed *bool}`, `DetailsProps{Open *bool}`, `DialogProps{Open *bool}`.

- [ ] **Step 1: Write the failing test**

```go
func TestBlockquoteAttrs(t *testing.T) {
	cite := "https://example.com"
	out := Render(Blockquote(BlockquoteProps{Cite: &cite}))
	if !strings.Contains(out, `cite="https://example.com"`) {
		t.Fatalf("got: %s", out)
	}
}

func TestDelInsAttrs(t *testing.T) {
	cite, datetime := "https://example.com", "2026-01-01"
	del := Render(Del(DelProps{Cite: &cite, Datetime: &datetime}))
	ins := Render(Ins(InsProps{Cite: &cite, Datetime: &datetime}))
	for _, out := range []string{del, ins} {
		if !strings.Contains(out, `cite="https://example.com"`) || !strings.Contains(out, `datetime="2026-01-01"`) {
			t.Fatalf("got: %s", out)
		}
	}
}

func TestQAttrs(t *testing.T) {
	cite := "https://example.com"
	out := Render(Q(QProps{Cite: &cite}))
	if !strings.Contains(out, `cite="https://example.com"`) {
		t.Fatalf("got: %s", out)
	}
}

func TestTimeAttrs(t *testing.T) {
	dt := "2026-01-01"
	out := Render(Time(TimeProps{Datetime: &dt}))
	if !strings.Contains(out, `datetime="2026-01-01"`) {
		t.Fatalf("got: %s", out)
	}
}

func TestDataAttrs(t *testing.T) {
	v := "42"
	out := Render(Data(DataProps{Value: &v}))
	if !strings.Contains(out, `value="42"`) {
		t.Fatalf("got: %s", out)
	}
}

func TestLiAttrs(t *testing.T) {
	v := "5"
	out := Render(Li(LiProps{Value: &v}))
	if !strings.Contains(out, `value="5"`) {
		t.Fatalf("got: %s", out)
	}
}

func TestOlAttrs(t *testing.T) {
	start, typ := "3", "A"
	reversed := true
	out := Render(Ol(OlProps{Start: &start, Type: &typ, Reversed: &reversed}))
	if !strings.Contains(out, `start="3"`) || !strings.Contains(out, `type="A"`) || !strings.Contains(out, "reversed") {
		t.Fatalf("got: %s", out)
	}
}

func TestDetailsDialogOpen(t *testing.T) {
	open := true
	d := Render(Details(DetailsProps{Open: &open}))
	dl := Render(Dialog(DialogProps{Open: &open}))
	for _, out := range []string{d, dl} {
		if !strings.Contains(out, "open") {
			t.Fatalf("got: %s", out)
		}
	}
}
```

- [ ] **Step 2: Run test to verify it fails**

Run: `go test -run 'TestBlockquote|TestDelIns|TestQAttrs|TestTimeAttrs|TestDataAttrs|TestLiAttrs|TestOlAttrs|TestDetailsDialog' ./...`
Expected: FAIL (unknown fields)

- [ ] **Step 3: Write minimal implementation**

`blockquote.go`:
```go
package html

type BlockquoteProps struct {
	ID, Class, Style, Title, Cite *string
	Attributes                    []Attribute
}

func Blockquote(p BlockquoteProps, children ...Node) Node {
	a := globalAttrs(p.ID, p.Class, p.Style, p.Title, p.Attributes)
	if p.Cite != nil {
		a = append(a, Attribute{"cite", *p.Cite, false})
	}
	return elementNode{name: "blockquote", props: a, void: false, children: children}
}
```

`del.go`:
```go
package html

type DelProps struct {
	ID, Class, Style, Title, Cite, Datetime *string
	Attributes                              []Attribute
}

func Del(p DelProps, children ...Node) Node {
	a := globalAttrs(p.ID, p.Class, p.Style, p.Title, p.Attributes)
	if p.Cite != nil {
		a = append(a, Attribute{"cite", *p.Cite, false})
	}
	if p.Datetime != nil {
		a = append(a, Attribute{"datetime", *p.Datetime, false})
	}
	return elementNode{name: "del", props: a, void: false, children: children}
}
```

`ins.go`:
```go
package html

type InsProps struct {
	ID, Class, Style, Title, Cite, Datetime *string
	Attributes                              []Attribute
}

func Ins(p InsProps, children ...Node) Node {
	a := globalAttrs(p.ID, p.Class, p.Style, p.Title, p.Attributes)
	if p.Cite != nil {
		a = append(a, Attribute{"cite", *p.Cite, false})
	}
	if p.Datetime != nil {
		a = append(a, Attribute{"datetime", *p.Datetime, false})
	}
	return elementNode{name: "ins", props: a, void: false, children: children}
}
```

`q.go`:
```go
package html

type QProps struct {
	ID, Class, Style, Title, Cite *string
	Attributes                    []Attribute
}

func Q(p QProps, children ...Node) Node {
	a := globalAttrs(p.ID, p.Class, p.Style, p.Title, p.Attributes)
	if p.Cite != nil {
		a = append(a, Attribute{"cite", *p.Cite, false})
	}
	return elementNode{name: "q", props: a, void: false, children: children}
}
```

`time.go`:
```go
package html

type TimeProps struct {
	ID, Class, Style, Title, Datetime *string
	Attributes                        []Attribute
}

func Time(p TimeProps, children ...Node) Node {
	a := globalAttrs(p.ID, p.Class, p.Style, p.Title, p.Attributes)
	if p.Datetime != nil {
		a = append(a, Attribute{"datetime", *p.Datetime, false})
	}
	return elementNode{name: "time", props: a, void: false, children: children}
}
```

`data.go`:
```go
package html

type DataProps struct {
	ID, Class, Style, Title, Value *string
	Attributes                     []Attribute
}

func Data(p DataProps, children ...Node) Node {
	a := globalAttrs(p.ID, p.Class, p.Style, p.Title, p.Attributes)
	if p.Value != nil {
		a = append(a, Attribute{"value", *p.Value, false})
	}
	return elementNode{name: "data", props: a, void: false, children: children}
}
```

`li.go`:
```go
package html

type LiProps struct {
	ID, Class, Style, Title, Value *string
	Attributes                     []Attribute
}

func Li(p LiProps, children ...Node) Node {
	a := globalAttrs(p.ID, p.Class, p.Style, p.Title, p.Attributes)
	if p.Value != nil {
		a = append(a, Attribute{"value", *p.Value, false})
	}
	return elementNode{name: "li", props: a, void: false, children: children}
}
```

`ol.go`:
```go
package html

type OlProps struct {
	ID, Class, Style, Title, Start, Type *string
	Reversed                             *bool
	Attributes                           []Attribute
}

func Ol(p OlProps, children ...Node) Node {
	a := globalAttrs(p.ID, p.Class, p.Style, p.Title, p.Attributes)
	if p.Start != nil {
		a = append(a, Attribute{"start", *p.Start, false})
	}
	if p.Type != nil {
		a = append(a, Attribute{"type", *p.Type, false})
	}
	if p.Reversed != nil {
		a = append(a, Attribute{"reversed", "true", true})
	}
	return elementNode{name: "ol", props: a, void: false, children: children}
}
```

`details.go`:
```go
package html

type DetailsProps struct {
	ID, Class, Style, Title *string
	Open                    *bool
	Attributes              []Attribute
}

func Details(p DetailsProps, children ...Node) Node {
	a := globalAttrs(p.ID, p.Class, p.Style, p.Title, p.Attributes)
	if p.Open != nil {
		a = append(a, Attribute{"open", "true", true})
	}
	return elementNode{name: "details", props: a, void: false, children: children}
}
```

`dialog.go`:
```go
package html

type DialogProps struct {
	ID, Class, Style, Title *string
	Open                    *bool
	Attributes              []Attribute
}

func Dialog(p DialogProps, children ...Node) Node {
	a := globalAttrs(p.ID, p.Class, p.Style, p.Title, p.Attributes)
	if p.Open != nil {
		a = append(a, Attribute{"open", "true", true})
	}
	return elementNode{name: "dialog", props: a, void: false, children: children}
}
```

- [ ] **Step 4: Run test to verify it passes**

Run: `go test ./...`
Expected: PASS

- [ ] **Step 5: Commit**

```bash
git add blockquote.go del.go ins.go q.go time.go data.go li.go ol.go details.go dialog.go element_attrs_test.go
git commit -m "extend inline semantic elements with real attrs"
```

---

## Task 5: Table cell/structure elements (td, th, colgroup, table)

**Files:**
- Modify: `td.go`, `th.go`, `colgroup.go`, `table.go`
- Test: append to `element_attrs_test.go`

**Interfaces:**
- Produces: `TdProps{Colspan, Rowspan, Headers}`, `ThProps{Colspan, Rowspan, Headers, Scope, Abbr}`, `ColgroupProps{Span}`, `TableProps{Bgcolor}`.

- [ ] **Step 1: Write the failing test**

```go
func TestTdAttrs(t *testing.T) {
	colspan, rowspan, headers := "2", "1", "h1 h2"
	out := Render(Td(TdProps{Colspan: &colspan, Rowspan: &rowspan, Headers: &headers}))
	for _, want := range []string{`colspan="2"`, `rowspan="1"`, `headers="h1 h2"`} {
		if !strings.Contains(out, want) {
			t.Fatalf("expected %s in %s", want, out)
		}
	}
}

func TestThAttrs(t *testing.T) {
	scope, abbr := "col", "Description"
	out := Render(Th(ThProps{Scope: &scope, Abbr: &abbr}))
	if !strings.Contains(out, `scope="col"`) || !strings.Contains(out, `abbr="Description"`) {
		t.Fatalf("got: %s", out)
	}
}

func TestColgroupAttrs(t *testing.T) {
	span := "3"
	out := Render(Colgroup(ColgroupProps{Span: &span}))
	if !strings.Contains(out, `span="3"`) {
		t.Fatalf("got: %s", out)
	}
}

func TestTableAttrs(t *testing.T) {
	bgcolor := "#eee"
	out := Render(Table(TableProps{Bgcolor: &bgcolor}))
	if !strings.Contains(out, `bgcolor="#eee"`) {
		t.Fatalf("got: %s", out)
	}
}
```

- [ ] **Step 2: Run test to verify it fails**

Run: `go test -run 'TestTdAttrs|TestThAttrs|TestColgroupAttrs|TestTableAttrs' ./...`
Expected: FAIL (unknown fields)

- [ ] **Step 3: Write minimal implementation**

`td.go`:
```go
package html

type TdProps struct {
	ID, Class, Style, Title, Colspan, Rowspan, Headers *string
	Attributes                                         []Attribute
}

func Td(p TdProps, children ...Node) Node {
	a := globalAttrs(p.ID, p.Class, p.Style, p.Title, p.Attributes)
	if p.Colspan != nil {
		a = append(a, Attribute{"colspan", *p.Colspan, false})
	}
	if p.Rowspan != nil {
		a = append(a, Attribute{"rowspan", *p.Rowspan, false})
	}
	if p.Headers != nil {
		a = append(a, Attribute{"headers", *p.Headers, false})
	}
	return elementNode{name: "td", props: a, void: false, children: children}
}
```

`th.go`:
```go
package html

type ThProps struct {
	ID, Class, Style, Title, Colspan, Rowspan, Headers, Scope, Abbr *string
	Attributes                                                      []Attribute
}

func Th(p ThProps, children ...Node) Node {
	a := globalAttrs(p.ID, p.Class, p.Style, p.Title, p.Attributes)
	if p.Colspan != nil {
		a = append(a, Attribute{"colspan", *p.Colspan, false})
	}
	if p.Rowspan != nil {
		a = append(a, Attribute{"rowspan", *p.Rowspan, false})
	}
	if p.Headers != nil {
		a = append(a, Attribute{"headers", *p.Headers, false})
	}
	if p.Scope != nil {
		a = append(a, Attribute{"scope", *p.Scope, false})
	}
	if p.Abbr != nil {
		a = append(a, Attribute{"abbr", *p.Abbr, false})
	}
	return elementNode{name: "th", props: a, void: false, children: children}
}
```

`colgroup.go`:
```go
package html

type ColgroupProps struct {
	ID, Class, Style, Title, Span *string
	Attributes                    []Attribute
}

func Colgroup(p ColgroupProps, children ...Node) Node {
	a := globalAttrs(p.ID, p.Class, p.Style, p.Title, p.Attributes)
	if p.Span != nil {
		a = append(a, Attribute{"span", *p.Span, false})
	}
	return elementNode{name: "colgroup", props: a, void: false, children: children}
}
```

`table.go`:
```go
package html

type TableProps struct {
	ID, Class, Style, Title, Bgcolor *string
	Attributes                       []Attribute
}

func Table(p TableProps, children ...Node) Node {
	a := globalAttrs(p.ID, p.Class, p.Style, p.Title, p.Attributes)
	if p.Bgcolor != nil {
		a = append(a, Attribute{"bgcolor", *p.Bgcolor, false})
	}
	return elementNode{name: "table", props: a, void: false, children: children}
}
```

- [ ] **Step 4: Run test to verify it passes**

Run: `go test ./...`
Expected: PASS

- [ ] **Step 5: Commit**

```bash
git add td.go th.go colgroup.go table.go element_attrs_test.go
git commit -m "extend table structure elements with real attrs"
```

---

## Task 6: Media/embedding elements (audio, video, canvas, map, object, iframe)

**Files:**
- Modify: `audio.go`, `video.go`, `canvas.go`, `map.go`, `object.go`, `iframe.go`
- Test: append to `element_attrs_test.go`

**Interfaces:**
- Produces: `AudioProps{Src, Preload, Controls, Autoplay, Loop, Muted *bool}`, `VideoProps{Src, Poster, Width, Height, Preload, Controls, Autoplay, Loop, Muted *bool}`, `CanvasProps{Width, Height}`, `MapProps{Name}`, `ObjectProps{Data, Type, Width, Height}`, `IframeProps{Src, Width, Height, Name, Allow, Loading}`.

- [ ] **Step 1: Write the failing test**

```go
func TestAudioAttrs(t *testing.T) {
	src, preload := "song.mp3", "auto"
	controls, autoplay, loop, muted := true, true, true, true
	out := Render(Audio(AudioProps{Src: &src, Preload: &preload, Controls: &controls, Autoplay: &autoplay, Loop: &loop, Muted: &muted}))
	for _, want := range []string{`src="song.mp3"`, `preload="auto"`, "controls", "autoplay", "loop", "muted"} {
		if !strings.Contains(out, want) {
			t.Fatalf("expected %s in %s", want, out)
		}
	}
}

func TestVideoAttrs(t *testing.T) {
	src, poster, w, h := "movie.mp4", "poster.jpg", "640", "480"
	controls := true
	out := Render(Video(VideoProps{Src: &src, Poster: &poster, Width: &w, Height: &h, Controls: &controls}))
	for _, want := range []string{`src="movie.mp4"`, `poster="poster.jpg"`, `width="640"`, `height="480"`, "controls"} {
		if !strings.Contains(out, want) {
			t.Fatalf("expected %s in %s", want, out)
		}
	}
}

func TestCanvasAttrs(t *testing.T) {
	w, h := "300", "150"
	out := Render(Canvas(CanvasProps{Width: &w, Height: &h}))
	if !strings.Contains(out, `width="300"`) || !strings.Contains(out, `height="150"`) {
		t.Fatalf("got: %s", out)
	}
}

func TestMapAttrs(t *testing.T) {
	name := "image-map"
	out := Render(Map(MapProps{Name: &name}))
	if !strings.Contains(out, `name="image-map"`) {
		t.Fatalf("got: %s", out)
	}
}

func TestObjectAttrs(t *testing.T) {
	data, typ, w, h := "movie.swf", "application/x-shockwave-flash", "400", "300"
	out := Render(Object(ObjectProps{Data: &data, Type: &typ, Width: &w, Height: &h}))
	for _, want := range []string{`data="movie.swf"`, `type="application/x-shockwave-flash"`, `width="400"`, `height="300"`} {
		if !strings.Contains(out, want) {
			t.Fatalf("expected %s in %s", want, out)
		}
	}
}

func TestIframeAttrs(t *testing.T) {
	src, w, h, name, allow, loading := "page.html", "600", "400", "frame1", "fullscreen", "lazy"
	out := Render(Iframe(IframeProps{Src: &src, Width: &w, Height: &h, Name: &name, Allow: &allow, Loading: &loading}))
	for _, want := range []string{`src="page.html"`, `width="600"`, `height="400"`, `name="frame1"`, `allow="fullscreen"`, `loading="lazy"`} {
		if !strings.Contains(out, want) {
			t.Fatalf("expected %s in %s", want, out)
		}
	}
}
```

- [ ] **Step 2: Run test to verify it fails**

Run: `go test -run 'TestAudioAttrs|TestVideoAttrs|TestCanvasAttrs|TestMapAttrs|TestObjectAttrs|TestIframeAttrs' ./...`
Expected: FAIL (unknown fields)

- [ ] **Step 3: Write minimal implementation**

`audio.go`:
```go
package html

type AudioProps struct {
	ID, Class, Style, Title, Src, Preload         *string
	Controls, Autoplay, Loop, Muted               *bool
	Attributes                                    []Attribute
}

func Audio(p AudioProps, children ...Node) Node {
	a := globalAttrs(p.ID, p.Class, p.Style, p.Title, p.Attributes)
	if p.Src != nil {
		a = append(a, Attribute{"src", *p.Src, false})
	}
	if p.Preload != nil {
		a = append(a, Attribute{"preload", *p.Preload, false})
	}
	if p.Controls != nil {
		a = append(a, Attribute{"controls", "true", true})
	}
	if p.Autoplay != nil {
		a = append(a, Attribute{"autoplay", "true", true})
	}
	if p.Loop != nil {
		a = append(a, Attribute{"loop", "true", true})
	}
	if p.Muted != nil {
		a = append(a, Attribute{"muted", "true", true})
	}
	return elementNode{name: "audio", props: a, void: false, children: children}
}
```

`video.go`:
```go
package html

type VideoProps struct {
	ID, Class, Style, Title, Src, Poster, Width, Height, Preload *string
	Controls, Autoplay, Loop, Muted                              *bool
	Attributes                                                   []Attribute
}

func Video(p VideoProps, children ...Node) Node {
	a := globalAttrs(p.ID, p.Class, p.Style, p.Title, p.Attributes)
	if p.Src != nil {
		a = append(a, Attribute{"src", *p.Src, false})
	}
	if p.Poster != nil {
		a = append(a, Attribute{"poster", *p.Poster, false})
	}
	if p.Width != nil {
		a = append(a, Attribute{"width", *p.Width, false})
	}
	if p.Height != nil {
		a = append(a, Attribute{"height", *p.Height, false})
	}
	if p.Preload != nil {
		a = append(a, Attribute{"preload", *p.Preload, false})
	}
	if p.Controls != nil {
		a = append(a, Attribute{"controls", "true", true})
	}
	if p.Autoplay != nil {
		a = append(a, Attribute{"autoplay", "true", true})
	}
	if p.Loop != nil {
		a = append(a, Attribute{"loop", "true", true})
	}
	if p.Muted != nil {
		a = append(a, Attribute{"muted", "true", true})
	}
	return elementNode{name: "video", props: a, void: false, children: children}
}
```

`canvas.go`:
```go
package html

type CanvasProps struct {
	ID, Class, Style, Title, Width, Height *string
	Attributes                             []Attribute
}

func Canvas(p CanvasProps, children ...Node) Node {
	a := globalAttrs(p.ID, p.Class, p.Style, p.Title, p.Attributes)
	if p.Width != nil {
		a = append(a, Attribute{"width", *p.Width, false})
	}
	if p.Height != nil {
		a = append(a, Attribute{"height", *p.Height, false})
	}
	return elementNode{name: "canvas", props: a, void: false, children: children}
}
```

`map.go`:
```go
package html

type MapProps struct {
	ID, Class, Style, Title, Name *string
	Attributes                    []Attribute
}

func Map(p MapProps, children ...Node) Node {
	a := globalAttrs(p.ID, p.Class, p.Style, p.Title, p.Attributes)
	if p.Name != nil {
		a = append(a, Attribute{"name", *p.Name, false})
	}
	return elementNode{name: "map", props: a, void: false, children: children}
}
```

`object.go`:
```go
package html

type ObjectProps struct {
	ID, Class, Style, Title, Data, Type, Width, Height *string
	Attributes                                         []Attribute
}

func Object(p ObjectProps, children ...Node) Node {
	a := globalAttrs(p.ID, p.Class, p.Style, p.Title, p.Attributes)
	if p.Data != nil {
		a = append(a, Attribute{"data", *p.Data, false})
	}
	if p.Type != nil {
		a = append(a, Attribute{"type", *p.Type, false})
	}
	if p.Width != nil {
		a = append(a, Attribute{"width", *p.Width, false})
	}
	if p.Height != nil {
		a = append(a, Attribute{"height", *p.Height, false})
	}
	return elementNode{name: "object", props: a, void: false, children: children}
}
```

`iframe.go`:
```go
package html

type IframeProps struct {
	ID, Class, Style, Title, Src, Width, Height, Name, Allow, Loading *string
	Attributes                                                       []Attribute
}

func Iframe(p IframeProps, children ...Node) Node {
	a := globalAttrs(p.ID, p.Class, p.Style, p.Title, p.Attributes)
	if p.Src != nil {
		a = append(a, Attribute{"src", *p.Src, false})
	}
	if p.Width != nil {
		a = append(a, Attribute{"width", *p.Width, false})
	}
	if p.Height != nil {
		a = append(a, Attribute{"height", *p.Height, false})
	}
	if p.Name != nil {
		a = append(a, Attribute{"name", *p.Name, false})
	}
	if p.Allow != nil {
		a = append(a, Attribute{"allow", *p.Allow, false})
	}
	if p.Loading != nil {
		a = append(a, Attribute{"loading", *p.Loading, false})
	}
	return elementNode{name: "iframe", props: a, void: false, children: children}
}
```

- [ ] **Step 4: Run test to verify it passes**

Run: `go test ./...`
Expected: PASS

- [ ] **Step 5: Commit**

```bash
git add audio.go video.go canvas.go map.go object.go iframe.go element_attrs_test.go
git commit -m "extend media/embedding elements with real attrs"
```

---

## Task 7: Form structure elements (form, label, fieldset, output, meter, progress)

**Files:**
- Modify: `form.go`, `label.go`, `fieldset.go`, `output.go`, `meter.go`, `progress.go`
- Test: append to `element_attrs_test.go`

**Interfaces:**
- Produces: `FormProps{Action, Method, Name, Target, Enctype, Novalidate *bool}`, `LabelProps{For}`, `FieldsetProps{Name, Disabled *bool}`, `OutputProps{For, Name}`, `MeterProps{Value, Min, Max, Low, High, Optimum}`, `ProgressProps{Value, Max}`.

- [ ] **Step 1: Write the failing test**

```go
func TestFormAttrs(t *testing.T) {
	action, method, name, target, enctype := "/submit", "post", "signup", "_top", "multipart/form-data"
	novalidate := true
	out := Render(Form(FormProps{Action: &action, Method: &method, Name: &name, Target: &target, Enctype: &enctype, Novalidate: &novalidate}))
	for _, want := range []string{`action="/submit"`, `method="post"`, `name="signup"`, `target="_top"`, `enctype="multipart/form-data"`, "novalidate"} {
		if !strings.Contains(out, want) {
			t.Fatalf("expected %s in %s", want, out)
		}
	}
}

func TestLabelAttrs(t *testing.T) {
	forID := "email"
	out := Render(Label(LabelProps{For: &forID}))
	if !strings.Contains(out, `for="email"`) {
		t.Fatalf("got: %s", out)
	}
}

func TestFieldsetAttrs(t *testing.T) {
	name := "personal"
	disabled := true
	out := Render(Fieldset(FieldsetProps{Name: &name, Disabled: &disabled}))
	if !strings.Contains(out, `name="personal"`) || !strings.Contains(out, "disabled") {
		t.Fatalf("got: %s", out)
	}
}

func TestOutputAttrs(t *testing.T) {
	forIDs, name := "a b", "result"
	out := Render(Output(OutputProps{For: &forIDs, Name: &name}))
	if !strings.Contains(out, `for="a b"`) || !strings.Contains(out, `name="result"`) {
		t.Fatalf("got: %s", out)
	}
}

func TestMeterAttrs(t *testing.T) {
	value, min, max, low, high, optimum := "6", "0", "10", "3", "8", "5"
	out := Render(Meter(MeterProps{Value: &value, Min: &min, Max: &max, Low: &low, High: &high, Optimum: &optimum}))
	for _, want := range []string{`value="6"`, `min="0"`, `max="10"`, `low="3"`, `high="8"`, `optimum="5"`} {
		if !strings.Contains(out, want) {
			t.Fatalf("expected %s in %s", want, out)
		}
	}
}

func TestProgressAttrs(t *testing.T) {
	value, max := "70", "100"
	out := Render(Progress(ProgressProps{Value: &value, Max: &max}))
	if !strings.Contains(out, `value="70"`) || !strings.Contains(out, `max="100"`) {
		t.Fatalf("got: %s", out)
	}
}
```

- [ ] **Step 2: Run test to verify it fails**

Run: `go test -run 'TestFormAttrs|TestLabelAttrs|TestFieldsetAttrs|TestOutputAttrs|TestMeterAttrs|TestProgressAttrs' ./...`
Expected: FAIL (unknown fields)

- [ ] **Step 3: Write minimal implementation**

`form.go`:
```go
package html

type FormProps struct {
	ID, Class, Style, Title, Action, Method, Name, Target, Enctype *string
	Novalidate                                                     *bool
	Attributes                                                     []Attribute
}

func Form(p FormProps, children ...Node) Node {
	a := globalAttrs(p.ID, p.Class, p.Style, p.Title, p.Attributes)
	if p.Action != nil {
		a = append(a, Attribute{"action", *p.Action, false})
	}
	if p.Method != nil {
		a = append(a, Attribute{"method", *p.Method, false})
	}
	if p.Name != nil {
		a = append(a, Attribute{"name", *p.Name, false})
	}
	if p.Target != nil {
		a = append(a, Attribute{"target", *p.Target, false})
	}
	if p.Enctype != nil {
		a = append(a, Attribute{"enctype", *p.Enctype, false})
	}
	if p.Novalidate != nil {
		a = append(a, Attribute{"novalidate", "true", true})
	}
	return elementNode{name: "form", props: a, void: false, children: children}
}
```

`label.go`:
```go
package html

type LabelProps struct {
	ID, Class, Style, Title, For *string
	Attributes                   []Attribute
}

func Label(p LabelProps, children ...Node) Node {
	a := globalAttrs(p.ID, p.Class, p.Style, p.Title, p.Attributes)
	if p.For != nil {
		a = append(a, Attribute{"for", *p.For, false})
	}
	return elementNode{name: "label", props: a, void: false, children: children}
}
```

`fieldset.go`:
```go
package html

type FieldsetProps struct {
	ID, Class, Style, Title, Name *string
	Disabled                      *bool
	Attributes                    []Attribute
}

func Fieldset(p FieldsetProps, children ...Node) Node {
	a := globalAttrs(p.ID, p.Class, p.Style, p.Title, p.Attributes)
	if p.Name != nil {
		a = append(a, Attribute{"name", *p.Name, false})
	}
	if p.Disabled != nil {
		a = append(a, Attribute{"disabled", "true", true})
	}
	return elementNode{name: "fieldset", props: a, void: false, children: children}
}
```

`output.go`:
```go
package html

type OutputProps struct {
	ID, Class, Style, Title, For, Name *string
	Attributes                         []Attribute
}

func Output(p OutputProps, children ...Node) Node {
	a := globalAttrs(p.ID, p.Class, p.Style, p.Title, p.Attributes)
	if p.For != nil {
		a = append(a, Attribute{"for", *p.For, false})
	}
	if p.Name != nil {
		a = append(a, Attribute{"name", *p.Name, false})
	}
	return elementNode{name: "output", props: a, void: false, children: children}
}
```

`meter.go`:
```go
package html

type MeterProps struct {
	ID, Class, Style, Title, Value, Min, Max, Low, High, Optimum *string
	Attributes                                                  []Attribute
}

func Meter(p MeterProps, children ...Node) Node {
	a := globalAttrs(p.ID, p.Class, p.Style, p.Title, p.Attributes)
	if p.Value != nil {
		a = append(a, Attribute{"value", *p.Value, false})
	}
	if p.Min != nil {
		a = append(a, Attribute{"min", *p.Min, false})
	}
	if p.Max != nil {
		a = append(a, Attribute{"max", *p.Max, false})
	}
	if p.Low != nil {
		a = append(a, Attribute{"low", *p.Low, false})
	}
	if p.High != nil {
		a = append(a, Attribute{"high", *p.High, false})
	}
	if p.Optimum != nil {
		a = append(a, Attribute{"optimum", *p.Optimum, false})
	}
	return elementNode{name: "meter", props: a, void: false, children: children}
}
```

`progress.go`:
```go
package html

type ProgressProps struct {
	ID, Class, Style, Title, Value, Max *string
	Attributes                          []Attribute
}

func Progress(p ProgressProps, children ...Node) Node {
	a := globalAttrs(p.ID, p.Class, p.Style, p.Title, p.Attributes)
	if p.Value != nil {
		a = append(a, Attribute{"value", *p.Value, false})
	}
	if p.Max != nil {
		a = append(a, Attribute{"max", *p.Max, false})
	}
	return elementNode{name: "progress", props: a, void: false, children: children}
}
```

- [ ] **Step 4: Run test to verify it passes**

Run: `go test ./...`
Expected: PASS

- [ ] **Step 5: Commit**

```bash
git add form.go label.go fieldset.go output.go meter.go progress.go element_attrs_test.go
git commit -m "extend form structure elements with real attrs"
```

---

## Task 8: Select/option elements (select, option, optgroup, textarea)

**Files:**
- Modify: `select.go`, `option.go`, `optgroup.go`, `textarea.go`
- Test: append to `element_attrs_test.go`

**Interfaces:**
- Produces: `SelectProps{Name, Size, Multiple, Required, Disabled *bool}`, `OptionProps{Value, Label, Selected, Disabled *bool}`, `OptgroupProps{Label, Disabled *bool}`, `TextareaProps{Name, Rows, Cols, Placeholder, Maxlength, Required, Disabled, Readonly *bool}`.

- [ ] **Step 1: Write the failing test**

```go
func TestSelectAttrs(t *testing.T) {
	name, size := "color", "1"
	multiple, required, disabled := true, true, true
	out := Render(Select(SelectProps{Name: &name, Size: &size, Multiple: &multiple, Required: &required, Disabled: &disabled}))
	for _, want := range []string{`name="color"`, `size="1"`, "multiple", "required", "disabled"} {
		if !strings.Contains(out, want) {
			t.Fatalf("expected %s in %s", want, out)
		}
	}
}

func TestOptionAttrs(t *testing.T) {
	value, label := "red", "Red"
	selected, disabled := true, true
	out := Render(Option(OptionProps{Value: &value, Label: &label, Selected: &selected, Disabled: &disabled}))
	for _, want := range []string{`value="red"`, `label="Red"`, "selected", "disabled"} {
		if !strings.Contains(out, want) {
			t.Fatalf("expected %s in %s", want, out)
		}
	}
}

func TestOptgroupAttrs(t *testing.T) {
	label := "Warm colors"
	disabled := true
	out := Render(Optgroup(OptgroupProps{Label: &label, Disabled: &disabled}))
	if !strings.Contains(out, `label="Warm colors"`) || !strings.Contains(out, "disabled") {
		t.Fatalf("got: %s", out)
	}
}

func TestTextareaAttrs(t *testing.T) {
	name, rows, cols, placeholder, maxlength := "bio", "4", "40", "Tell us about yourself", "200"
	required, disabled, readonly := true, true, true
	out := Render(Textarea(TextareaProps{Name: &name, Rows: &rows, Cols: &cols, Placeholder: &placeholder, Maxlength: &maxlength, Required: &required, Disabled: &disabled, Readonly: &readonly}))
	for _, want := range []string{`name="bio"`, `rows="4"`, `cols="40"`, `placeholder="Tell us about yourself"`, `maxlength="200"`, "required", "disabled", "readonly"} {
		if !strings.Contains(out, want) {
			t.Fatalf("expected %s in %s", want, out)
		}
	}
}
```

- [ ] **Step 2: Run test to verify it fails**

Run: `go test -run 'TestSelectAttrs|TestOptionAttrs|TestOptgroupAttrs|TestTextareaAttrs' ./...`
Expected: FAIL (unknown fields)

- [ ] **Step 3: Write minimal implementation**

`select.go`:
```go
package html

type SelectProps struct {
	ID, Class, Style, Title, Name, Size *string
	Multiple, Required, Disabled        *bool
	Attributes                          []Attribute
}

func Select(p SelectProps, children ...Node) Node {
	a := globalAttrs(p.ID, p.Class, p.Style, p.Title, p.Attributes)
	if p.Name != nil {
		a = append(a, Attribute{"name", *p.Name, false})
	}
	if p.Size != nil {
		a = append(a, Attribute{"size", *p.Size, false})
	}
	if p.Multiple != nil {
		a = append(a, Attribute{"multiple", "true", true})
	}
	if p.Required != nil {
		a = append(a, Attribute{"required", "true", true})
	}
	if p.Disabled != nil {
		a = append(a, Attribute{"disabled", "true", true})
	}
	return elementNode{name: "select", props: a, void: false, children: children}
}
```

`option.go`:
```go
package html

type OptionProps struct {
	ID, Class, Style, Title, Value, Label *string
	Selected, Disabled                    *bool
	Attributes                            []Attribute
}

func Option(p OptionProps, children ...Node) Node {
	a := globalAttrs(p.ID, p.Class, p.Style, p.Title, p.Attributes)
	if p.Value != nil {
		a = append(a, Attribute{"value", *p.Value, false})
	}
	if p.Label != nil {
		a = append(a, Attribute{"label", *p.Label, false})
	}
	if p.Selected != nil {
		a = append(a, Attribute{"selected", "true", true})
	}
	if p.Disabled != nil {
		a = append(a, Attribute{"disabled", "true", true})
	}
	return elementNode{name: "option", props: a, void: false, children: children}
}
```

`optgroup.go`:
```go
package html

type OptgroupProps struct {
	ID, Class, Style, Title, Label *string
	Disabled                       *bool
	Attributes                     []Attribute
}

func Optgroup(p OptgroupProps, children ...Node) Node {
	a := globalAttrs(p.ID, p.Class, p.Style, p.Title, p.Attributes)
	if p.Label != nil {
		a = append(a, Attribute{"label", *p.Label, false})
	}
	if p.Disabled != nil {
		a = append(a, Attribute{"disabled", "true", true})
	}
	return elementNode{name: "optgroup", props: a, void: false, children: children}
}
```

`textarea.go`:
```go
package html

type TextareaProps struct {
	ID, Class, Style, Title, Name, Rows, Cols, Placeholder, Maxlength *string
	Required, Disabled, Readonly                                     *bool
	Attributes                                                       []Attribute
}

func Textarea(p TextareaProps, children ...Node) Node {
	a := globalAttrs(p.ID, p.Class, p.Style, p.Title, p.Attributes)
	if p.Name != nil {
		a = append(a, Attribute{"name", *p.Name, false})
	}
	if p.Rows != nil {
		a = append(a, Attribute{"rows", *p.Rows, false})
	}
	if p.Cols != nil {
		a = append(a, Attribute{"cols", *p.Cols, false})
	}
	if p.Placeholder != nil {
		a = append(a, Attribute{"placeholder", *p.Placeholder, false})
	}
	if p.Maxlength != nil {
		a = append(a, Attribute{"maxlength", *p.Maxlength, false})
	}
	if p.Required != nil {
		a = append(a, Attribute{"required", "true", true})
	}
	if p.Disabled != nil {
		a = append(a, Attribute{"disabled", "true", true})
	}
	if p.Readonly != nil {
		a = append(a, Attribute{"readonly", "true", true})
	}
	return elementNode{name: "textarea", props: a, void: false, children: children}
}
```

- [ ] **Step 4: Run test to verify it passes**

Run: `go test ./...`
Expected: PASS

- [ ] **Step 5: Commit**

```bash
git add select.go option.go optgroup.go textarea.go element_attrs_test.go
git commit -m "extend select/option/textarea elements with real attrs"
```

---

## Task 9: Document metadata elements (script_tag, style, body)

**Files:**
- Modify: `script_tag.go`, `style.go`, `body.go`
- Test: append to `element_attrs_test.go`

**Interfaces:**
- Produces: `ScriptTagProps{Src, Type, Async, Defer *bool}`, `StyleProps{Media, Type}`, `BodyProps{Background, Bgcolor, Text, Link, Vlink, Alink}`.

- [ ] **Step 1: Write the failing test**

```go
func TestScriptTagAttrs(t *testing.T) {
	src, typ := "app.js", "module"
	async, defer_ := true, true
	out := Render(ScriptTag(ScriptTagProps{Src: &src, Type: &typ, Async: &async, Defer: &defer_}))
	for _, want := range []string{`src="app.js"`, `type="module"`, "async", "defer"} {
		if !strings.Contains(out, want) {
			t.Fatalf("expected %s in %s", want, out)
		}
	}
}

func TestStyleAttrs(t *testing.T) {
	media, typ := "print", "text/css"
	out := Render(Style(StyleProps{Media: &media, Type: &typ}))
	if !strings.Contains(out, `media="print"`) || !strings.Contains(out, `type="text/css"`) {
		t.Fatalf("got: %s", out)
	}
}

func TestBodyAttrs(t *testing.T) {
	bg, bgcolor, text, link, vlink, alink := "bg.png", "#fff", "#000", "#00f", "#808", "#f00"
	out := Render(Body(BodyProps{Background: &bg, Bgcolor: &bgcolor, Text: &text, Link: &link, Vlink: &vlink, Alink: &alink}))
	for _, want := range []string{`background="bg.png"`, `bgcolor="#fff"`, `text="#000"`, `link="#00f"`, `vlink="#808"`, `alink="#f00"`} {
		if !strings.Contains(out, want) {
			t.Fatalf("expected %s in %s", want, out)
		}
	}
}
```

- [ ] **Step 2: Run test to verify it fails**

Run: `go test -run 'TestScriptTagAttrs|TestStyleAttrs|TestBodyAttrs' ./...`
Expected: FAIL (unknown fields)

- [ ] **Step 3: Write minimal implementation**

`script_tag.go`:
```go
package html

type ScriptTagProps struct {
	ID, Class, Style, Title, Src, Type *string
	Async, Defer                       *bool
	Attributes                         []Attribute
}

func ScriptTag(p ScriptTagProps, children ...Node) Node {
	a := globalAttrs(p.ID, p.Class, p.Style, p.Title, p.Attributes)
	if p.Src != nil {
		a = append(a, Attribute{"src", *p.Src, false})
	}
	if p.Type != nil {
		a = append(a, Attribute{"type", *p.Type, false})
	}
	if p.Async != nil {
		a = append(a, Attribute{"async", "true", true})
	}
	if p.Defer != nil {
		a = append(a, Attribute{"defer", "true", true})
	}
	return elementNode{name: "script", props: a, void: false, children: children}
}
```

`style.go`:
```go
package html

type StyleProps struct {
	ID, Class, Style, Title, Media, Type *string
	Attributes                           []Attribute
}

func Style(p StyleProps, children ...Node) Node {
	a := globalAttrs(p.ID, p.Class, p.Style, p.Title, p.Attributes)
	if p.Media != nil {
		a = append(a, Attribute{"media", *p.Media, false})
	}
	if p.Type != nil {
		a = append(a, Attribute{"type", *p.Type, false})
	}
	return elementNode{name: "style", props: a, void: false, children: children}
}
```

`body.go`:
```go
package html

type BodyProps struct {
	ID, Class, Style, Title, Background, Bgcolor, Text, Link, Vlink, Alink *string
	Attributes                                                             []Attribute
}

func Body(p BodyProps, children ...Node) Node {
	a := globalAttrs(p.ID, p.Class, p.Style, p.Title, p.Attributes)
	if p.Background != nil {
		a = append(a, Attribute{"background", *p.Background, false})
	}
	if p.Bgcolor != nil {
		a = append(a, Attribute{"bgcolor", *p.Bgcolor, false})
	}
	if p.Text != nil {
		a = append(a, Attribute{"text", *p.Text, false})
	}
	if p.Link != nil {
		a = append(a, Attribute{"link", *p.Link, false})
	}
	if p.Vlink != nil {
		a = append(a, Attribute{"vlink", *p.Vlink, false})
	}
	if p.Alink != nil {
		a = append(a, Attribute{"alink", *p.Alink, false})
	}
	return elementNode{name: "body", props: a, void: false, children: children}
}
```

- [ ] **Step 4: Run test to verify it passes**

Run: `go test ./...`
Expected: PASS

- [ ] **Step 5: Commit**

```bash
git add script_tag.go style.go body.go element_attrs_test.go
git commit -m "extend script/style/body elements with real attrs"
```

---

## Task 10: Deprecated legacy elements (font, marquee, frame, frameset, dir)

**Files:**
- Modify: `font.go`, `marquee.go`, `frame.go`, `frameset.go`, `dir.go`
- Test: append to `element_attrs_test.go`

**Interfaces:**
- Produces: `FontProps{Size, Color, Face}`, `MarqueeProps{Behavior, Direction, Scrollamount, Scrolldelay, Loop}`, `FrameProps{Src, Name, Noresize *bool}`, `FramesetProps{Rows, Cols}`, `DirProps{Compact *bool}`.

- [ ] **Step 1: Write the failing test**

```go
func TestFontAttrs(t *testing.T) {
	size, color, face := "4", "red", "Arial"
	out := Render(Font(FontProps{Size: &size, Color: &color, Face: &face}))
	for _, want := range []string{`size="4"`, `color="red"`, `face="Arial"`} {
		if !strings.Contains(out, want) {
			t.Fatalf("expected %s in %s", want, out)
		}
	}
}

func TestMarqueeAttrs(t *testing.T) {
	behavior, direction, amount, delay, loop := "scroll", "left", "6", "85", "3"
	out := Render(Marquee(MarqueeProps{Behavior: &behavior, Direction: &direction, Scrollamount: &amount, Scrolldelay: &delay, Loop: &loop}))
	for _, want := range []string{`behavior="scroll"`, `direction="left"`, `scrollamount="6"`, `scrolldelay="85"`, `loop="3"`} {
		if !strings.Contains(out, want) {
			t.Fatalf("expected %s in %s", want, out)
		}
	}
}

func TestFrameAttrs(t *testing.T) {
	src, name := "nav.html", "navframe"
	noresize := true
	out := Render(Frame(FrameProps{Src: &src, Name: &name, Noresize: &noresize}))
	if !strings.Contains(out, `src="nav.html"`) || !strings.Contains(out, `name="navframe"`) || !strings.Contains(out, "noresize") {
		t.Fatalf("got: %s", out)
	}
}

func TestFramesetAttrs(t *testing.T) {
	rows, cols := "50%,50%", "*"
	out := Render(Frameset(FramesetProps{Rows: &rows, Cols: &cols}))
	if !strings.Contains(out, `rows="50%,50%"`) || !strings.Contains(out, `cols="*"`) {
		t.Fatalf("got: %s", out)
	}
}

func TestDirAttrs(t *testing.T) {
	compact := true
	out := Render(Dir(DirProps{Compact: &compact}))
	if !strings.Contains(out, "compact") {
		t.Fatalf("got: %s", out)
	}
}
```

- [ ] **Step 2: Run test to verify it fails**

Run: `go test -run 'TestFontAttrs|TestMarqueeAttrs|TestFrameAttrs|TestFramesetAttrs|TestDirAttrs' ./...`
Expected: FAIL (unknown fields)

- [ ] **Step 3: Write minimal implementation**

`font.go`:
```go
package html

type FontProps struct {
	ID, Class, Style, Title, Size, Color, Face *string
	Attributes                                 []Attribute
}

func Font(p FontProps, children ...Node) Node {
	a := globalAttrs(p.ID, p.Class, p.Style, p.Title, p.Attributes)
	if p.Size != nil {
		a = append(a, Attribute{"size", *p.Size, false})
	}
	if p.Color != nil {
		a = append(a, Attribute{"color", *p.Color, false})
	}
	if p.Face != nil {
		a = append(a, Attribute{"face", *p.Face, false})
	}
	return elementNode{name: "font", props: a, void: false, children: children}
}
```

`marquee.go`:
```go
package html

type MarqueeProps struct {
	ID, Class, Style, Title, Behavior, Direction, Scrollamount, Scrolldelay, Loop *string
	Attributes                                                                   []Attribute
}

func Marquee(p MarqueeProps, children ...Node) Node {
	a := globalAttrs(p.ID, p.Class, p.Style, p.Title, p.Attributes)
	if p.Behavior != nil {
		a = append(a, Attribute{"behavior", *p.Behavior, false})
	}
	if p.Direction != nil {
		a = append(a, Attribute{"direction", *p.Direction, false})
	}
	if p.Scrollamount != nil {
		a = append(a, Attribute{"scrollamount", *p.Scrollamount, false})
	}
	if p.Scrolldelay != nil {
		a = append(a, Attribute{"scrolldelay", *p.Scrolldelay, false})
	}
	if p.Loop != nil {
		a = append(a, Attribute{"loop", *p.Loop, false})
	}
	return elementNode{name: "marquee", props: a, void: false, children: children}
}
```

`frame.go`:
```go
package html

type FrameProps struct {
	ID, Class, Style, Title, Src, Name *string
	Noresize                           *bool
	Attributes                         []Attribute
}

func Frame(p FrameProps, children ...Node) Node {
	a := globalAttrs(p.ID, p.Class, p.Style, p.Title, p.Attributes)
	if p.Src != nil {
		a = append(a, Attribute{"src", *p.Src, false})
	}
	if p.Name != nil {
		a = append(a, Attribute{"name", *p.Name, false})
	}
	if p.Noresize != nil {
		a = append(a, Attribute{"noresize", "true", true})
	}
	return elementNode{name: "frame", props: a, void: false, children: children}
}
```

`frameset.go`:
```go
package html

type FramesetProps struct {
	ID, Class, Style, Title, Rows, Cols *string
	Attributes                          []Attribute
}

func Frameset(p FramesetProps, children ...Node) Node {
	a := globalAttrs(p.ID, p.Class, p.Style, p.Title, p.Attributes)
	if p.Rows != nil {
		a = append(a, Attribute{"rows", *p.Rows, false})
	}
	if p.Cols != nil {
		a = append(a, Attribute{"cols", *p.Cols, false})
	}
	return elementNode{name: "frameset", props: a, void: false, children: children}
}
```

`dir.go`:
```go
package html

type DirProps struct {
	ID, Class, Style, Title *string
	Compact                 *bool
	Attributes              []Attribute
}

func Dir(p DirProps, children ...Node) Node {
	a := globalAttrs(p.ID, p.Class, p.Style, p.Title, p.Attributes)
	if p.Compact != nil {
		a = append(a, Attribute{"compact", "true", true})
	}
	return elementNode{name: "dir", props: a, void: false, children: children}
}
```

- [ ] **Step 4: Run test to verify it passes**

Run: `go test ./...`
Expected: PASS

- [ ] **Step 5: Commit**

```bash
git add font.go marquee.go frame.go frameset.go dir.go element_attrs_test.go
git commit -m "extend deprecated legacy elements with real attrs"
```

---

## Task 11: Delete README Limitations, add Scoped styles/scripts section, update notes

**Files:**
- Modify: `README.md`, `notes/Architecture.md`, `notes/Decisions Log.md`, `notes/Codebase Overview.md`

**Interfaces:**
- Consumes: `ScopedStyle` (Task 1), `ScopedScript` (Task 2), and every `*Props` field added in Tasks 3-10.

- [ ] **Step 1: Verify every element with real attributes now exposes them**

Run:
```bash
grep -L '"href"\|"src"\|"name"\|"value"\|"span"\|"cite"' *.go 2>/dev/null | grep -v _test.go
```
Expected: this lists only files genuinely having no element-specific attributes (text-level elements like `b.go`, `i.go`, `strong.go`, container-only elements like `div.go`, `ul.go`, structural no-op elements like `thead.go`, and the pipeline files). Manually confirm no file in Tasks 3-10's file lists appears — if one does, that task's edit didn't land; re-check it before continuing.

- [ ] **Step 2: Replace README's Limitations section with a Scoped styles and scripts section**

Remove:
```markdown
## Limitations

- Most element files only expose `ID`, `Class`, `Style`, `Title`,
  and `Attributes` — reach other attributes through `Attributes` until a
  given file is hand-extended (as `input.go`, `button.go`, `a.go`, `meta.go`,
  and `link.go` already are).
- `Fragment(children...)` groups nodes into a `[]Node` for composing pieces
  together; there's still no component abstraction beyond plain `Node`.
- No CSS/JS scoping — `<style>`/`<script>` content is emitted as-is; keeping
  it collision-free across a page is your responsibility.
```

Add, in its place (same position in the file, right before `## License`):

```markdown
## Scoped styles and scripts

`ScopedStyle` namespaces a CSS string's top-level selectors under one id,
so two components' styles can share a page without colliding:

```sh
cat <<'GO' > /tmp/scoped_example.go
package main

import (
	"fmt"

	html "github.com/nsatyasrikar/html"
)

func main() {
	fmt.Println(html.Render(html.ScopedStyle("card-1", ".title { color: red; }")))
}
GO
go run /tmp/scoped_example.go
```

```
<style> #card-1 .title { color: red; }</style>
```

`ScopedScript` wraps a script body in an IIFE, so every `var`/`let`/
`const`/`function` it declares stays out of the global scope:

```go
html.Render(html.ScopedScript("var count = 0;"))
// <script>(function(){
// var count = 0;
// })();</script>
```
```

(Verify both `go run` outputs actually match before committing — this repo's README convention requires every example be checked against real output, not assumed.)

- [ ] **Step 3: Update `notes/Architecture.md`'s `## Elements` section**

Replace the sentence "Most files still only expose the generic four props; a handful gained element-specific fields..." with: "Every element with real HTML attributes now exposes them as typed fields (see the full list added across `docs/superpowers/plans/2026-09-09-close-limitations.md`); elements with no attributes beyond the generic four (most inline text-level elements, pure containers) are unchanged."

- [ ] **Step 4: Update `notes/Decisions Log.md`**

Add a new entry after `## htmlgen removal`:

```markdown
## Closing the Limitations list

Followed up the htmlgen removal by hand-writing every element's real
attributes (see `docs/superpowers/specs/2026-09-09-close-limitations-design.md`
and its plan) rather than leaving them behind `Attributes`. Also added
`ScopedStyle` (selector-prefix string rewrite) and `ScopedScript` (IIFE
wrapper) so CSS/JS collision risk is closed without a CSS/JS parser —
IIFE scoping is a complete guarantee via JS's own semantics, not a
best-effort string hack. README's `## Limitations` section is gone.
```

- [ ] **Step 5: Update `notes/Codebase Overview.md`'s Quick facts**

Change "Every element file (`div.go`, `p.go`, ...) follows one template..." bullet to also link `[[Decisions Log#Closing the Limitations list]]` instead of `#htmlgen removal`.

- [ ] **Step 6: Full verification**

Run: `go build ./... && go vet ./... && go test ./... && gofmt -l .`
Expected: build/vet/test all pass; `gofmt -l .` prints nothing (no unformatted files).

- [ ] **Step 7: Commit**

```bash
git add README.md notes/Architecture.md "notes/Decisions Log.md" "notes/Codebase Overview.md"
git commit -m "delete README Limitations, document ScopedStyle/ScopedScript"
```
