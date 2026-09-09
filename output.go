package html

type OutputProps struct {
	ID, Class, Style, Title *string
	Attributes              []Attribute
}

func Output(p OutputProps, children ...Node) Node {
	a := globalAttrs(p.ID, p.Class, p.Style, p.Title, p.Attributes)
	return elementNode{name: "output", props: a, void: false, children: children}
}
