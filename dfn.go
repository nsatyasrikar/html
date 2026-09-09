package html

type DfnProps struct {
	ID, Class, Style, Title *string
	Attributes              []Attribute
}

func Dfn(p DfnProps, children ...Node) Node {
	a := globalAttrs(p.ID, p.Class, p.Style, p.Title, p.Attributes)
	return elementNode{name: "dfn", props: a, void: false, children: children}
}
