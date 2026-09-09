package html

type FormProps struct {
	ID, Class, Style, Title, Action, Method, Name, Target, Enctype *string
	Novalidate                                                     *bool
	Attributes                                                     []Attribute
}

func Form(p FormProps, children ...Node) Node {
	a := globalAttrs(p.ID, p.Class, p.Style, p.Title, p.Attributes)
	if p.Action != nil {
		a = append(a, Attribute{"action", *p.Action, false})
	}
	if p.Method != nil {
		a = append(a, Attribute{"method", *p.Method, false})
	}
	if p.Name != nil {
		a = append(a, Attribute{"name", *p.Name, false})
	}
	if p.Target != nil {
		a = append(a, Attribute{"target", *p.Target, false})
	}
	if p.Enctype != nil {
		a = append(a, Attribute{"enctype", *p.Enctype, false})
	}
	if p.Novalidate != nil {
		a = append(a, Attribute{"novalidate", "true", true})
	}
	return elementNode{name: "form", props: a, void: false, children: children}
}
