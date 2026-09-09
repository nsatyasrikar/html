package html

type StyleProps struct {
	ID, Class, Style, Title, Media, Type *string
	Attributes                           []Attribute
}

func Style(p StyleProps, children ...Node) Node {
	a := globalAttrs(p.ID, p.Class, p.Style, p.Title, p.Attributes)
	if p.Media != nil {
		a = append(a, Attribute{"media", *p.Media, false})
	}
	if p.Type != nil {
		a = append(a, Attribute{"type", *p.Type, false})
	}
	return elementNode{name: "style", props: a, void: false, children: children}
}
