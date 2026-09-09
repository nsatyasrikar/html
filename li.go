package html

type LiProps struct {
	ID, Class, Style, Title, Value *string
	Attributes                     []Attribute
}

func Li(p LiProps, children ...Node) Node {
	a := globalAttrs(p.ID, p.Class, p.Style, p.Title, p.Attributes)
	if p.Value != nil {
		a = append(a, Attribute{"value", *p.Value, false})
	}
	return elementNode{name: "li", props: a, void: false, children: children}
}
