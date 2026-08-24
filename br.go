package html

type BrProps struct {
	ID, Class, Style, Title *string
	Attributes              []Attribute
}

func Br(p BrProps) Node {
	return elementNode{name: "br", props: globalAttrs(p.ID, p.Class, p.Style, p.Title, p.Attributes), void: true}
}
