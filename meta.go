package html

type MetaProps struct {
	ID, Class, Style, Title, Name, Content *string
	Attributes                            []Attribute
}

func Meta(p MetaProps) Node {
	a := globalAttrs(p.ID, p.Class, p.Style, p.Title, p.Attributes)
	if p.Name != nil {
		a = append(a, Attribute{"name", *p.Name, false})
	}
	if p.Content != nil {
		a = append(a, Attribute{"content", *p.Content, false})
	}
	return elementNode{name: "meta", props: a, void: true}
}
