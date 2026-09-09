package html

type OptionProps struct {
	ID, Class, Style, Title, Value, Label *string
	Selected                              *bool
	Disabled                              *bool
	Attributes                            []Attribute
}

func Option(p OptionProps, children ...Node) Node {
	a := globalAttrs(p.ID, p.Class, p.Style, p.Title, p.Attributes)
	if p.Value != nil {
		a = append(a, Attribute{"value", *p.Value, false})
	}
	if p.Label != nil {
		a = append(a, Attribute{"label", *p.Label, false})
	}
	if p.Selected != nil {
		a = append(a, Attribute{"selected", "true", true})
	}
	if p.Disabled != nil {
		a = append(a, Attribute{"disabled", "true", true})
	}
	return elementNode{name: "option", props: a, void: false, children: children}
}
