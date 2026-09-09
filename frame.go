package html

type FrameProps struct {
	ID, Class, Style, Title, Src, Name *string
	Noresize                           *bool
	Attributes                         []Attribute
}

func Frame(p FrameProps, children ...Node) Node {
	a := globalAttrs(p.ID, p.Class, p.Style, p.Title, p.Attributes)
	if p.Src != nil {
		a = append(a, Attribute{"src", *p.Src, false})
	}
	if p.Name != nil {
		a = append(a, Attribute{"name", *p.Name, false})
	}
	if p.Noresize != nil {
		a = append(a, Attribute{"noresize", "true", true})
	}
	return elementNode{name: "frame", props: a, void: false, children: children}
}
