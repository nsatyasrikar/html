package html

import "testing"

func TestScriptCodeEscapesClosingTag(t *testing.T) {
	got := Render(ScriptTag(ScriptTagProps{}, ScriptCode("A</ScRiPt>B")))
	want := "<script>A<\\/ScRiPt>B</script>"
	if got != want {
		t.Fatalf("got %q want %q", got, want)
	}
}
