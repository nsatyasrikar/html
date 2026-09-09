package html

type UlProps struct {
	ID, Class, Style, Title *string
	Attributes              []Attribute
}

func Ul(p UlProps, children ...Node) Node {
	a := globalAttrs(p.ID, p.Class, p.Style, p.Title, p.Attributes)
	return elementNode{name: "ul", props: a, void: false, children: children}
}
