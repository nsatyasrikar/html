package html

type DatalistProps struct {
	ID, Class, Style, Title *string
	Attributes              []Attribute
}

func Datalist(p DatalistProps, children ...Node) Node {
	a := globalAttrs(p.ID, p.Class, p.Style, p.Title, p.Attributes)
	return elementNode{name: "datalist", props: a, void: false, children: children}
}
