package html

type LinkProps struct {
	ID, Class, Style, Title, Rel, Href *string
	Attributes                        []Attribute
}

func Link(p LinkProps) Node {
	a := globalAttrs(p.ID, p.Class, p.Style, p.Title, p.Attributes)
	if p.Rel != nil {
		a = append(a, Attribute{"rel", *p.Rel, false})
	}
	if p.Href != nil {
		a = append(a, Attribute{"href", *p.Href, false})
	}
	return elementNode{name: "link", props: a, void: true}
}
