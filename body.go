package html

type BodyProps struct {
	ID, Class, Style, Title, Background, Bgcolor, Text, Link, Vlink, Alink *string
	Attributes                                                             []Attribute
}

func Body(p BodyProps, children ...Node) Node {
	a := globalAttrs(p.ID, p.Class, p.Style, p.Title, p.Attributes)
	if p.Background != nil {
		a = append(a, Attribute{"background", *p.Background, false})
	}
	if p.Bgcolor != nil {
		a = append(a, Attribute{"bgcolor", *p.Bgcolor, false})
	}
	if p.Text != nil {
		a = append(a, Attribute{"text", *p.Text, false})
	}
	if p.Link != nil {
		a = append(a, Attribute{"link", *p.Link, false})
	}
	if p.Vlink != nil {
		a = append(a, Attribute{"vlink", *p.Vlink, false})
	}
	if p.Alink != nil {
		a = append(a, Attribute{"alink", *p.Alink, false})
	}
	return elementNode{name: "body", props: a, void: false, children: children}
}
