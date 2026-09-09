package html

type CodeProps struct {
	ID, Class, Style, Title *string
	Attributes              []Attribute
}

func Code(p CodeProps, children ...Node) Node {
	a := globalAttrs(p.ID, p.Class, p.Style, p.Title, p.Attributes)
	return elementNode{name: "code", props: a, void: false, children: children}
}
