package html

type AreaProps struct {
	ID, Class, Style, Title *string
	Attributes              []Attribute
}

func Area(p AreaProps) Node {
	a := globalAttrs(p.ID, p.Class, p.Style, p.Title, p.Attributes)
	return elementNode{name: "area", props: a, void: true}
}
