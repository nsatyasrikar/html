package html

type SelectedcontentProps struct {
	ID, Class, Style, Title *string
	Attributes              []Attribute
}

func Selectedcontent(p SelectedcontentProps, children ...Node) Node {
	a := globalAttrs(p.ID, p.Class, p.Style, p.Title, p.Attributes)
	return elementNode{name: "selectedcontent", props: a, void: false, children: children}
}
