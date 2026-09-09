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
