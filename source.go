package html

type SourceProps struct {
	ID, Class, Style, Title *string
	Attributes              []Attribute
}

func Source(p SourceProps) Node {
	a := globalAttrs(p.ID, p.Class, p.Style, p.Title, p.Attributes)
	return elementNode{name: "source", props: a, void: true}
}
