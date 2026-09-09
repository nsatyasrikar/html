package html

type TimeProps struct {
	ID, Class, Style, Title, Datetime *string
	Attributes                        []Attribute
}

func Time(p TimeProps, children ...Node) Node {
	a := globalAttrs(p.ID, p.Class, p.Style, p.Title, p.Attributes)
	if p.Datetime != nil {
		a = append(a, Attribute{"datetime", *p.Datetime, false})
	}
	return elementNode{name: "time", props: a, void: false, children: children}
}
