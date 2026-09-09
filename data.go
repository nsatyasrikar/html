package html

type DataProps struct {
	ID, Class, Style, Title, Value *string
	Attributes                     []Attribute
}

func Data(p DataProps, children ...Node) Node {
	a := globalAttrs(p.ID, p.Class, p.Style, p.Title, p.Attributes)
	if p.Value != nil {
		a = append(a, Attribute{"value", *p.Value, false})
	}
	return elementNode{name: "data", props: a, void: false, children: children}
}
