package html

type H5Props struct {
	ID, Class, Style, Title *string
	Attributes              []Attribute
}

func H5(p H5Props, children ...Node) Node {
	a := globalAttrs(p.ID, p.Class, p.Style, p.Title, p.Attributes)
	return elementNode{name: "h5", props: a, void: false, children: children}
}
