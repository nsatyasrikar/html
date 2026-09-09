package html

type ImgProps struct {
	ID, Class, Style, Title, Src, Alt, Width, Height, Srcset, Sizes, Loading *string
	Attributes                                                               []Attribute
}

func Img(p ImgProps) Node {
	a := globalAttrs(p.ID, p.Class, p.Style, p.Title, p.Attributes)
	if p.Src != nil {
		a = append(a, Attribute{"src", *p.Src, false})
	}
	if p.Alt != nil {
		a = append(a, Attribute{"alt", *p.Alt, false})
	}
	if p.Width != nil {
		a = append(a, Attribute{"width", *p.Width, false})
	}
	if p.Height != nil {
		a = append(a, Attribute{"height", *p.Height, false})
	}
	if p.Srcset != nil {
		a = append(a, Attribute{"srcset", *p.Srcset, false})
	}
	if p.Sizes != nil {
		a = append(a, Attribute{"sizes", *p.Sizes, false})
	}
	if p.Loading != nil {
		a = append(a, Attribute{"loading", *p.Loading, false})
	}
	return elementNode{name: "img", props: a, void: true}
}
