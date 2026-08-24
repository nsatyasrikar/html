package html

type globalProps struct {
	ID, Class, Style, Title *string
	Attributes              []Attribute
}

func globalAttrs(id, class, style, title *string, extra []Attribute) []Attribute {
	a := make([]Attribute, 0, len(extra)+4)
	if id != nil {
		a = append(a, Attribute{"id", *id, false})
	}
	if class != nil {
		a = append(a, Attribute{"class", *class, false})
	}
	if style != nil {
		a = append(a, Attribute{"style", *style, false})
	}
	if title != nil {
		a = append(a, Attribute{"title", *title, false})
	}
	return append(a, customAttributes(extra)...)
}
