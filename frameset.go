package html

type FramesetProps struct {
	ID, Class, Style, Title, Rows, Cols *string
	Attributes                          []Attribute
}

func Frameset(p FramesetProps, children ...Node) Node {
	a := globalAttrs(p.ID, p.Class, p.Style, p.Title, p.Attributes)
	if p.Rows != nil {
		a = append(a, Attribute{"rows", *p.Rows, false})
	}
	if p.Cols != nil {
		a = append(a, Attribute{"cols", *p.Cols, false})
	}
	return elementNode{name: "frameset", props: a, void: false, children: children}
}
