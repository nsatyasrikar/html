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
