package html

type TrProps struct {
	ID, Class, Style, Title *string
	Attributes              []Attribute
}

func Tr(p TrProps, children ...Node) Node {
	a := globalAttrs(p.ID, p.Class, p.Style, p.Title, p.Attributes)
	return elementNode{name: "tr", props: a, void: false, children: children}
}
