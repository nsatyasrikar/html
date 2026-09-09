package html

type ThProps struct {
	ID, Class, Style, Title, Colspan, Rowspan, Headers, Scope, Abbr *string
	Attributes                                                      []Attribute
}

func Th(p ThProps, children ...Node) Node {
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
	if p.Scope != nil {
		a = append(a, Attribute{"scope", *p.Scope, false})
	}
	if p.Abbr != nil {
		a = append(a, Attribute{"abbr", *p.Abbr, false})
	}
	return elementNode{name: "th", props: a, void: false, children: children}
}
