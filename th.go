package html

type ThProps struct {
	ID, Class, Style, Title *string
	Attributes              []Attribute
}

func Th(p ThProps, children ...Node) Node {
	a := globalAttrs(p.ID, p.Class, p.Style, p.Title, p.Attributes)
	return elementNode{name: "th", props: a, void: false, children: children}
}
