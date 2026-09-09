package html

type FontProps struct {
	ID, Class, Style, Title, Size, Color, Face *string
	Attributes                                 []Attribute
}

func Font(p FontProps, children ...Node) Node {
	a := globalAttrs(p.ID, p.Class, p.Style, p.Title, p.Attributes)
	if p.Size != nil {
		a = append(a, Attribute{"size", *p.Size, false})
	}
	if p.Color != nil {
		a = append(a, Attribute{"color", *p.Color, false})
	}
	if p.Face != nil {
		a = append(a, Attribute{"face", *p.Face, false})
	}
	return elementNode{name: "font", props: a, void: false, children: children}
}
