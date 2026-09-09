package html

type SourceProps struct {
	ID, Class, Style, Title, Src, Type, Srcset, Sizes, Media *string
	Attributes                                               []Attribute
}

func Source(p SourceProps) Node {
	a := globalAttrs(p.ID, p.Class, p.Style, p.Title, p.Attributes)
	if p.Src != nil {
		a = append(a, Attribute{"src", *p.Src, false})
	}
	if p.Type != nil {
		a = append(a, Attribute{"type", *p.Type, false})
	}
	if p.Srcset != nil {
		a = append(a, Attribute{"srcset", *p.Srcset, false})
	}
	if p.Sizes != nil {
		a = append(a, Attribute{"sizes", *p.Sizes, false})
	}
	if p.Media != nil {
		a = append(a, Attribute{"media", *p.Media, false})
	}
	return elementNode{name: "source", props: a, void: true}
}
