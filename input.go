package html

type InputType string

const (
	InputTypeText     InputType = "text"
	InputTypeEmail    InputType = "email"
	InputTypePassword InputType = "password"
	InputTypeNumber   InputType = "number"
)

type InputProps struct {
	ID, Class, Style, Title, Name, Value, Placeholder *string
	Type                                              *InputType
	Disabled, Required                                *bool
	OnInput, OnChange                                 *Script
	Attributes                                        []Attribute
}

func Input(p InputProps) Node {
	a := globalAttrs(p.ID, p.Class, p.Style, p.Title, p.Attributes)
	if p.Name != nil {
		a = append(a, Attribute{"name", *p.Name, false})
	}
	if p.Value != nil {
		a = append(a, Attribute{"value", *p.Value, false})
	}
	if p.Placeholder != nil {
		a = append(a, Attribute{"placeholder", *p.Placeholder, false})
	}
	if p.Type != nil {
		a = append(a, Attribute{"type", string(*p.Type), false})
	}
	if p.Disabled != nil {
		a = append(a, Attribute{"disabled", "true", true})
	}
	if p.Required != nil {
		a = append(a, Attribute{"required", "true", true})
	}
	if p.OnInput != nil {
		a = append(a, Attribute{"oninput", string(*p.OnInput), false})
	}
	if p.OnChange != nil {
		a = append(a, Attribute{"onchange", string(*p.OnChange), false})
	}
	return elementNode{name: "input", props: a, void: true}
}
