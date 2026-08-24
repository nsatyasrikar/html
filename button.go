package html

type ButtonProps struct {
	ID, Class, Style, Title, Name, Value *string
	Disabled                             *bool
	OnClick                              *Script
	Attributes                           []Attribute
}

func Button(p ButtonProps, children ...Node) Node {
	a := globalAttrs(p.ID, p.Class, p.Style, p.Title, p.Attributes)
	if p.Name != nil {
		a = append(a, Attribute{"name", *p.Name, false})
	}
	if p.Value != nil {
		a = append(a, Attribute{"value", *p.Value, false})
	}
	if p.Disabled != nil {
		a = append(a, Attribute{"disabled", "true", true})
	}
	if p.OnClick != nil {
		a = append(a, Attribute{"onclick", string(*p.OnClick), false})
	}
	return elementNode{name: "button", props: a, children: children}
}
