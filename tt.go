package html

type TtProps struct {
	ID, Class, Style, Title *string
	Attributes              []Attribute
}

func Tt(p TtProps, children ...Node) Node {
	a := globalAttrs(p.ID, p.Class, p.Style, p.Title, p.Attributes)
	return elementNode{name: "tt", props: a, void: false, children: children}
}
