package html

type DetailsProps struct {
	ID, Class, Style, Title *string
	Open                    *bool
	Attributes              []Attribute
}

func Details(p DetailsProps, children ...Node) Node {
	a := globalAttrs(p.ID, p.Class, p.Style, p.Title, p.Attributes)
	if p.Open != nil {
		a = append(a, Attribute{"open", "true", true})
	}
	return elementNode{name: "details", props: a, void: false, children: children}
}
