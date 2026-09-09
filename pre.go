package html

type PreProps struct {
	ID, Class, Style, Title *string
	Attributes              []Attribute
}

func Pre(p PreProps, children ...Node) Node {
	a := globalAttrs(p.ID, p.Class, p.Style, p.Title, p.Attributes)
	return elementNode{name: "pre", props: a, void: false, children: children}
}
