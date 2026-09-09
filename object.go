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
