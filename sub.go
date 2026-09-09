package html

type SubProps struct {
	ID, Class, Style, Title *string
	Attributes              []Attribute
}

func Sub(p SubProps, children ...Node) Node {
	a := globalAttrs(p.ID, p.Class, p.Style, p.Title, p.Attributes)
	return elementNode{name: "sub", props: a, void: false, children: children}
}
