package html

import (
	"strings"
	"testing"
)

func TestScopedScriptWrapsInIIFE(t *testing.T) {
	out := Render(ScopedScript("var x = 1;"))
	if !strings.Contains(out, "(function(){") || !strings.Contains(out, "})();") {
		t.Fatal(out)
	}
}
func TestScopedScriptBodyNotBareGlobal(t *testing.T) {
	out := Render(ScopedScript("var leaked = 1;"))
	closeIdx := strings.Index(out, "})();")
	if closeIdx < 0 || strings.Contains(out[closeIdx:], "var leaked") {
		t.Fatal(out)
	}
}
func TestScopedScriptRendersScriptTag(t *testing.T) {
	out := Render(ScopedScript("var x = 1;"))
	if !strings.HasPrefix(out, "<script>") || !strings.HasSuffix(out, "</script>") {
		t.Fatal(out)
	}
}
