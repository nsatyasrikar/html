package html

type MapProps struct {
	ID, Class, Style, Title, Name *string
	Attributes                    []Attribute
}

func Map(p MapProps, children ...Node) Node {
	a := globalAttrs(p.ID, p.Class, p.Style, p.Title, p.Attributes)
	if p.Name != nil {
		a = append(a, Attribute{"name", *p.Name, false})
	}
	return elementNode{name: "map", props: a, void: false, children: children}
}
