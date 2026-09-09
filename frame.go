package html

type FrameProps struct {
	ID, Class, Style, Title *string
	Attributes              []Attribute
}

func Frame(p FrameProps, children ...Node) Node {
	a := globalAttrs(p.ID, p.Class, p.Style, p.Title, p.Attributes)
	return elementNode{name: "frame", props: a, void: false, children: children}
}
