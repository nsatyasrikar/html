package html

type FieldsetProps struct {
	ID, Class, Style, Title, Name *string
	Disabled                      *bool
	Attributes                    []Attribute
}

func Fieldset(p FieldsetProps, children ...Node) Node {
	a := globalAttrs(p.ID, p.Class, p.Style, p.Title, p.Attributes)
	if p.Name != nil {
		a = append(a, Attribute{"name", *p.Name, false})
	}
	return elementNode{name: "fieldset", props: a, void: false, children: children}
}
