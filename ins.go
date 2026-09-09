package html

type InsProps struct {
	ID, Class, Style, Title, Cite, Datetime *string
	Attributes                              []Attribute
}

func Ins(p InsProps, children ...Node) Node {
	a := globalAttrs(p.ID, p.Class, p.Style, p.Title, p.Attributes)
	if p.Cite != nil {
		a = append(a, Attribute{"cite", *p.Cite, false})
	}
	if p.Datetime != nil {
		a = append(a, Attribute{"datetime", *p.Datetime, false})
	}
	return elementNode{name: "ins", props: a, void: false, children: children}
}
