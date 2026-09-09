package html

type QProps struct {
	ID, Class, Style, Title *string
	Attributes              []Attribute
}

func Q(p QProps, children ...Node) Node {
	a := globalAttrs(p.ID, p.Class, p.Style, p.Title, p.Attributes)
	return elementNode{name: "q", props: a, void: false, children: children}
}
