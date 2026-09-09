package html

type InsProps struct {
	ID, Class, Style, Title *string
	Attributes              []Attribute
}

func Ins(p InsProps, children ...Node) Node {
	a := globalAttrs(p.ID, p.Class, p.Style, p.Title, p.Attributes)
	return elementNode{name: "ins", props: a, void: false, children: children}
}
