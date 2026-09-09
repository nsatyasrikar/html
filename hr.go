package html

type HrProps struct {
	ID, Class, Style, Title, Color, Size, Width *string
	Attributes                                  []Attribute
}

func Hr(p HrProps) Node {
	a := globalAttrs(p.ID, p.Class, p.Style, p.Title, p.Attributes)
	if p.Color != nil {
		a = append(a, Attribute{"color", *p.Color, false})
	}
	if p.Size != nil {
		a = append(a, Attribute{"size", *p.Size, false})
	}
	if p.Width != nil {
		a = append(a, Attribute{"width", *p.Width, false})
	}
	return elementNode{name: "hr", props: a, void: true}
}
