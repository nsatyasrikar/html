package html

type ImgProps struct {
	ID, Class, Style, Title, Src, Alt *string
	Attributes                        []Attribute
}

func Img(p ImgProps) Node {
	a := globalAttrs(p.ID, p.Class, p.Style, p.Title, p.Attributes)
	if p.Src != nil {
		a = append(a, Attribute{"src", *p.Src, false})
	}
	if p.Alt != nil {
		a = append(a, Attribute{"alt", *p.Alt, false})
	}
	return elementNode{name: "img", props: a, void: true}
}
