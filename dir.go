package html

type DirProps struct {
	ID, Class, Style, Title *string
	Compact                 *bool
	Attributes              []Attribute
}

func Dir(p DirProps, children ...Node) Node {
	a := globalAttrs(p.ID, p.Class, p.Style, p.Title, p.Attributes)
	if p.Compact != nil {
		a = append(a, Attribute{"compact", "true", true})
	}
	return elementNode{name: "dir", props: a, void: false, children: children}
}
