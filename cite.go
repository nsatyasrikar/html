package html

type CiteProps struct {
	ID, Class, Style, Title *string
	Attributes              []Attribute
}

func Cite(p CiteProps, children ...Node) Node {
	a := globalAttrs(p.ID, p.Class, p.Style, p.Title, p.Attributes)
	return elementNode{name: "cite", props: a, void: false, children: children}
}
