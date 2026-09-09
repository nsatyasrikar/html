package html

type BaseProps struct {
	ID, Class, Style, Title, Href, Target *string
	Attributes                            []Attribute
}

func Base(p BaseProps) Node {
	a := globalAttrs(p.ID, p.Class, p.Style, p.Title, p.Attributes)
	if p.Href != nil {
		a = append(a, Attribute{"href", *p.Href, false})
	}
	if p.Target != nil {
		a = append(a, Attribute{"target", *p.Target, false})
	}
	return elementNode{name: "base", props: a, void: true}
}
