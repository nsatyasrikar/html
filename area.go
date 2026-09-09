package html

type AreaProps struct {
	ID, Class, Style, Title, Href, Alt, Coords, Shape, Target *string
	Attributes                                                []Attribute
}

func Area(p AreaProps) Node {
	a := globalAttrs(p.ID, p.Class, p.Style, p.Title, p.Attributes)
	if p.Href != nil {
		a = append(a, Attribute{"href", *p.Href, false})
	}
	if p.Alt != nil {
		a = append(a, Attribute{"alt", *p.Alt, false})
	}
	if p.Coords != nil {
		a = append(a, Attribute{"coords", *p.Coords, false})
	}
	if p.Shape != nil {
		a = append(a, Attribute{"shape", *p.Shape, false})
	}
	if p.Target != nil {
		a = append(a, Attribute{"target", *p.Target, false})
	}
	return elementNode{name: "area", props: a, void: true}
}
