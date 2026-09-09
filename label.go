package html

type LabelProps struct {
	ID, Class, Style, Title, For *string
	Attributes                   []Attribute
}

func Label(p LabelProps, children ...Node) Node {
	a := globalAttrs(p.ID, p.Class, p.Style, p.Title, p.Attributes)
	if p.For != nil {
		a = append(a, Attribute{"for", *p.For, false})
	}
	return elementNode{name: "label", props: a, void: false, children: children}
}
