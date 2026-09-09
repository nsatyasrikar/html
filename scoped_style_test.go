package html

import (
	"strings"
	"testing"
)

func TestScopedStyleSimpleSelector(t *testing.T) {
	out := Render(ScopedStyle("card-1", ".title { color: red; }"))
	if !strings.Contains(out, "#card-1 .title") {
		t.Fatal(out)
	}
}
func TestScopedStyleCommaSeparatedSelectors(t *testing.T) {
	out := Render(ScopedStyle("card-1", ".a, .b { color: red; }"))
	if !strings.Contains(out, "#card-1 .a") || !strings.Contains(out, "#card-1 .b") {
		t.Fatal(out)
	}
}
func TestScopedStyleAtRulePassthrough(t *testing.T) {
	out := Render(ScopedStyle("card-1", "@media (min-width: 600px) { .a { color: red; } }"))
	if !strings.Contains(out, "@media (min-width: 600px)") {
		t.Fatal(out)
	}
}
func TestScopedStyleRendersStyleTag(t *testing.T) {
	out := Render(ScopedStyle("card-1", ".a { color: red; }"))
	if !strings.HasPrefix(out, "<style>") || !strings.HasSuffix(out, "</style>") {
		t.Fatal(out)
	}
}
