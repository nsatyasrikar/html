package html

type ObjectProps struct {
	ID, Class, Style, Title *string
	Attributes              []Attribute
}

func Object(p ObjectProps, children ...Node) Node {
	a := globalAttrs(p.ID, p.Class, p.Style, p.Title, p.Attributes)
	return elementNode{name: "object", props: a, void: false, children: children}
}
