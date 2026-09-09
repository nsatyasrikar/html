package html

type OptgroupProps struct {
	ID, Class, Style, Title, Label *string
	Disabled                       *bool
	Attributes                     []Attribute
}

func Optgroup(p OptgroupProps, children ...Node) Node {
	a := globalAttrs(p.ID, p.Class, p.Style, p.Title, p.Attributes)
	if p.Label != nil {
		a = append(a, Attribute{"label", *p.Label, false})
	}
	if p.Disabled != nil {
		a = append(a, Attribute{"disabled", "true", true})
	}
	return elementNode{name: "optgroup", props: a, void: false, children: children}
}
