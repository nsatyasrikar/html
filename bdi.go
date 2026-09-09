package html

type BdiProps struct {
	ID, Class, Style, Title *string
	Attributes              []Attribute
}

func Bdi(p BdiProps, children ...Node) Node {
	a := globalAttrs(p.ID, p.Class, p.Style, p.Title, p.Attributes)
	return elementNode{name: "bdi", props: a, void: false, children: children}
}
