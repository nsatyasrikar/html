package html

type AudioProps struct {
	ID, Class, Style, Title *string
	Attributes              []Attribute
}

func Audio(p AudioProps, children ...Node) Node {
	a := globalAttrs(p.ID, p.Class, p.Style, p.Title, p.Attributes)
	return elementNode{name: "audio", props: a, void: false, children: children}
}
