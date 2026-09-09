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
