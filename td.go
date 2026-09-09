package html

type TdProps struct {
	ID, Class, Style, Title, Colspan, Rowspan, Headers *string
	Attributes                                         []Attribute
}

func Td(p TdProps, children ...Node) Node {
	a := globalAttrs(p.ID, p.Class, p.Style, p.Title, p.Attributes)
	if p.Colspan != nil {
		a = append(a, Attribute{"colspan", *p.Colspan, false})
	}
	if p.Rowspan != nil {
		a = append(a, Attribute{"rowspan", *p.Rowspan, false})
	}
	if p.Headers != nil {
		a = append(a, Attribute{"headers", *p.Headers, false})
	}
	return elementNode{name: "td", props: a, void: false, children: children}
}
