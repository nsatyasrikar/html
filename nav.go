package html

type NavProps struct {
	ID, Class, Style, Title *string
	Attributes              []Attribute
}

func Nav(p NavProps, children ...Node) Node {
	a := globalAttrs(p.ID, p.Class, p.Style, p.Title, p.Attributes)
	return elementNode{name: "nav", props: a, void: false, children: children}
}
