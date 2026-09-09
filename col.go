package html

type ColProps struct {
	ID, Class, Style, Title *string
	Attributes              []Attribute
}

func Col(p ColProps) Node {
	a := globalAttrs(p.ID, p.Class, p.Style, p.Title, p.Attributes)
	return elementNode{name: "col", props: a, void: true}
}
