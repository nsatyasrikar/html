package html

import "testing"

func TestRenderTree(t *testing.T) {
	got := Render(H1(H1Props{Class: Prop("heading")}, Text("Hello <world>")))
	want := `<h1 class="heading">Hello &lt;world&gt;</h1>`
	if got != want {
		t.Fatalf("got %q want %q", got, want)
	}
}
func TestVoidTagHasNoChildrenAPI(t *testing.T) {
	got := Render(Input(InputProps{Type: Prop(InputTypeEmail), Disabled: Prop(true)}))
	want := `<input type="email" disabled>`
	if got != want {
		t.Fatalf("got %q want %q", got, want)
	}
}
