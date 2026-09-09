package html

type AProps struct {
	ID, Class, Style, Title, Href *string
	Attributes                    []Attribute
}

func A(p AProps, children ...Node) Node {
	a := globalAttrs(p.ID, p.Class, p.Style, p.Title, p.Attributes)
	if p.Href != nil {
		a = append(a, Attribute{"href", *p.Href, false})
	}
	return elementNode{name: "a", props: a, void: false, children: children}
}
