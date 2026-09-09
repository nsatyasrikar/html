package html

type DelProps struct {
	ID, Class, Style, Title *string
	Attributes              []Attribute
}

func Del(p DelProps, children ...Node) Node {
	a := globalAttrs(p.ID, p.Class, p.Style, p.Title, p.Attributes)
	return elementNode{name: "del", props: a, void: false, children: children}
}
