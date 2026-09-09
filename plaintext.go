package html

type PlaintextProps struct {
	ID, Class, Style, Title *string
	Attributes              []Attribute
}

func Plaintext(p PlaintextProps, children ...Node) Node {
	a := globalAttrs(p.ID, p.Class, p.Style, p.Title, p.Attributes)
	return elementNode{name: "plaintext", props: a, void: false, children: children}
}
