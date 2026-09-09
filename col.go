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
