package html

type OlProps struct {
	ID, Class, Style, Title, Start, Type *string
	Reversed                             *bool
	Attributes                           []Attribute
}

func Ol(p OlProps, children ...Node) Node {
	a := globalAttrs(p.ID, p.Class, p.Style, p.Title, p.Attributes)
	if p.Start != nil {
		a = append(a, Attribute{"start", *p.Start, false})
	}
	if p.Type != nil {
		a = append(a, Attribute{"type", *p.Type, false})
	}
	if p.Reversed != nil {
		a = append(a, Attribute{"reversed", "true", true})
	}
	return elementNode{name: "ol", props: a, void: false, children: children}
}
