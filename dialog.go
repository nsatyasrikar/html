package html

type DialogProps struct {
	ID, Class, Style, Title *string
	Open                    *bool
	Attributes              []Attribute
}

func Dialog(p DialogProps, children ...Node) Node {
	a := globalAttrs(p.ID, p.Class, p.Style, p.Title, p.Attributes)
	if p.Open != nil {
		a = append(a, Attribute{"open", "true", true})
	}
	return elementNode{name: "dialog", props: a, void: false, children: children}
}
