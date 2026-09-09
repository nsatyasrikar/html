package html

type TableProps struct {
	ID, Class, Style, Title, Bgcolor *string
	Attributes                       []Attribute
}

func Table(p TableProps, children ...Node) Node {
	a := globalAttrs(p.ID, p.Class, p.Style, p.Title, p.Attributes)
	if p.Bgcolor != nil {
		a = append(a, Attribute{"bgcolor", *p.Bgcolor, false})
	}
	return elementNode{name: "table", props: a, void: false, children: children}
}
