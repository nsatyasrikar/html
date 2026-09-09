package html

type StyleProps struct {
	ID, Class, Style, Title *string
	Attributes              []Attribute
}

func Style(p StyleProps, children ...Node) Node {
	a := globalAttrs(p.ID, p.Class, p.Style, p.Title, p.Attributes)
	return elementNode{name: "style", props: a, void: false, children: children}
}
