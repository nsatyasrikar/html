package html

type TbodyProps struct {
	ID, Class, Style, Title *string
	Attributes              []Attribute
}

func Tbody(p TbodyProps, children ...Node) Node {
	a := globalAttrs(p.ID, p.Class, p.Style, p.Title, p.Attributes)
	return elementNode{name: "tbody", props: a, void: false, children: children}
}
