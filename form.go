package html

type FormProps struct {
	ID, Class, Style, Title *string
	Attributes              []Attribute
}

func Form(p FormProps, children ...Node) Node {
	a := globalAttrs(p.ID, p.Class, p.Style, p.Title, p.Attributes)
	return elementNode{name: "form", props: a, void: false, children: children}
}
