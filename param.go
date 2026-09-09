package html

type ParamProps struct {
	ID, Class, Style, Title, Name, Value *string
	Attributes                           []Attribute
}

func Param(p ParamProps) Node {
	a := globalAttrs(p.ID, p.Class, p.Style, p.Title, p.Attributes)
	if p.Name != nil {
		a = append(a, Attribute{"name", *p.Name, false})
	}
	if p.Value != nil {
		a = append(a, Attribute{"value", *p.Value, false})
	}
	return elementNode{name: "param", props: a, void: true}
}
