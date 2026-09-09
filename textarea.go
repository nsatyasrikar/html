package html

type TextareaProps struct {
	ID, Class, Style, Title, Name, Rows, Cols, Placeholder, Maxlength *string
	Required                                                          *bool
	Disabled                                                          *bool
	Readonly                                                          *bool
	Attributes                                                        []Attribute
}

func Textarea(p TextareaProps, children ...Node) Node {
	a := globalAttrs(p.ID, p.Class, p.Style, p.Title, p.Attributes)
	if p.Name != nil {
		a = append(a, Attribute{"name", *p.Name, false})
	}
	if p.Rows != nil {
		a = append(a, Attribute{"rows", *p.Rows, false})
	}
	if p.Cols != nil {
		a = append(a, Attribute{"cols", *p.Cols, false})
	}
	if p.Placeholder != nil {
		a = append(a, Attribute{"placeholder", *p.Placeholder, false})
	}
	if p.Maxlength != nil {
		a = append(a, Attribute{"maxlength", *p.Maxlength, false})
	}
	if p.Required != nil {
		a = append(a, Attribute{"required", "true", true})
	}
	if p.Disabled != nil {
		a = append(a, Attribute{"disabled", "true", true})
	}
	if p.Readonly != nil {
		a = append(a, Attribute{"readonly", "true", true})
	}
	return elementNode{name: "textarea", props: a, void: false, children: children}
}
