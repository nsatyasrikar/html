package html

type TitleProps struct {
	ID, Class, Style, Title *string
	Attributes              []Attribute
}

func Title(p TitleProps, children ...Node) Node {
	a := globalAttrs(p.ID, p.Class, p.Style, p.Title, p.Attributes)
	return elementNode{name: "title", props: a, void: false, children: children}
}
