package html

type DetailsProps struct {
	ID, Class, Style, Title *string
	Attributes              []Attribute
}

func Details(p DetailsProps, children ...Node) Node {
	a := globalAttrs(p.ID, p.Class, p.Style, p.Title, p.Attributes)
	return elementNode{name: "details", props: a, void: false, children: children}
}
