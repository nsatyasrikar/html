package html

type DivProps struct {
	ID, Class, Style, Title *string
	Attributes              []Attribute
}

func Div(p DivProps, children ...Node) Node {
	return elementNode{name: "div", props: globalAttrs(p.ID, p.Class, p.Style, p.Title, p.Attributes), children: children}
}
