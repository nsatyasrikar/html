package html

type ProgressProps struct {
	ID, Class, Style, Title, Value, Max *string
	Attributes                          []Attribute
}

func Progress(p ProgressProps, children ...Node) Node {
	a := globalAttrs(p.ID, p.Class, p.Style, p.Title, p.Attributes)
	if p.Value != nil {
		a = append(a, Attribute{"value", *p.Value, false})
	}
	if p.Max != nil {
		a = append(a, Attribute{"max", *p.Max, false})
	}
	return elementNode{name: "progress", props: a, void: false, children: children}
}
