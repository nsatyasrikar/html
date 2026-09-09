package html

type OutputProps struct {
	ID, Class, Style, Title, For, Name *string
	Attributes                         []Attribute
}

func Output(p OutputProps, children ...Node) Node {
	a := globalAttrs(p.ID, p.Class, p.Style, p.Title, p.Attributes)
	if p.For != nil {
		a = append(a, Attribute{"for", *p.For, false})
	}
	if p.Name != nil {
		a = append(a, Attribute{"name", *p.Name, false})
	}
	return elementNode{name: "output", props: a, void: false, children: children}
}
