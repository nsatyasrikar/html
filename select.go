package html

type SelectProps struct {
	ID, Class, Style, Title, Name, Size *string
	Multiple                            *bool
	Required                            *bool
	Disabled                            *bool
	Attributes                          []Attribute
}

func Select(p SelectProps, children ...Node) Node {
	a := globalAttrs(p.ID, p.Class, p.Style, p.Title, p.Attributes)
	if p.Name != nil {
		a = append(a, Attribute{"name", *p.Name, false})
	}
	if p.Size != nil {
		a = append(a, Attribute{"size", *p.Size, false})
	}
	if p.Multiple != nil {
		a = append(a, Attribute{"multiple", "true", true})
	}
	if p.Required != nil {
		a = append(a, Attribute{"required", "true", true})
	}
	if p.Disabled != nil {
		a = append(a, Attribute{"disabled", "true", true})
	}
	return elementNode{name: "select", props: a, void: false, children: children}
}
