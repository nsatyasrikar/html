package html

type QProps struct {
	ID, Class, Style, Title, Cite *string
	Attributes                    []Attribute
}

func Q(p QProps, children ...Node) Node {
	a := globalAttrs(p.ID, p.Class, p.Style, p.Title, p.Attributes)
	if p.Cite != nil {
		a = append(a, Attribute{"cite", *p.Cite, false})
	}
	return elementNode{name: "q", props: a, void: false, children: children}
}
