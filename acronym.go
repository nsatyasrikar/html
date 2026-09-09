package html

type AcronymProps struct {
	ID, Class, Style, Title *string
	Attributes              []Attribute
}

func Acronym(p AcronymProps, children ...Node) Node {
	a := globalAttrs(p.ID, p.Class, p.Style, p.Title, p.Attributes)
	return elementNode{name: "acronym", props: a, void: false, children: children}
}
