package html

type IProps struct {
	ID, Class, Style, Title *string
	Attributes              []Attribute
}

func I(p IProps, children ...Node) Node {
	a := globalAttrs(p.ID, p.Class, p.Style, p.Title, p.Attributes)
	return elementNode{name: "i", props: a, void: false, children: children}
}
