package html

type DtProps struct {
	ID, Class, Style, Title *string
	Attributes              []Attribute
}

func Dt(p DtProps, children ...Node) Node {
	a := globalAttrs(p.ID, p.Class, p.Style, p.Title, p.Attributes)
	return elementNode{name: "dt", props: a, void: false, children: children}
}
