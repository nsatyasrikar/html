package html

type TdProps struct {
	ID, Class, Style, Title *string
	Attributes              []Attribute
}

func Td(p TdProps, children ...Node) Node {
	a := globalAttrs(p.ID, p.Class, p.Style, p.Title, p.Attributes)
	return elementNode{name: "td", props: a, void: false, children: children}
}
