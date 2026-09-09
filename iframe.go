package html

type IframeProps struct {
	ID, Class, Style, Title, Src, Width, Height, Name, Allow, Loading *string
	Attributes                                                        []Attribute
}

func Iframe(p IframeProps, children ...Node) Node {
	a := globalAttrs(p.ID, p.Class, p.Style, p.Title, p.Attributes)
	if p.Src != nil {
		a = append(a, Attribute{"src", *p.Src, false})
	}
	if p.Width != nil {
		a = append(a, Attribute{"width", *p.Width, false})
	}
	if p.Height != nil {
		a = append(a, Attribute{"height", *p.Height, false})
	}
	if p.Name != nil {
		a = append(a, Attribute{"name", *p.Name, false})
	}
	if p.Allow != nil {
		a = append(a, Attribute{"allow", *p.Allow, false})
	}
	if p.Loading != nil {
		a = append(a, Attribute{"loading", *p.Loading, false})
	}
	return elementNode{name: "iframe", props: a, void: false, children: children}
}
