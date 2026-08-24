package html

type H1Props struct {
	ID, Class, Style, Title *string
	Attributes              []Attribute
}

func H1(p H1Props, children ...Node) Node {
	return elementNode{name: "h1", props: globalAttrs(p.ID, p.Class, p.Style, p.Title, p.Attributes), children: children}
}
