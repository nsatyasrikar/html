package html

type ScriptTagProps struct {
	ID, Class, Style, Title *string
	Attributes              []Attribute
}

func ScriptTag(p ScriptTagProps, children ...Node) Node {
	a := globalAttrs(p.ID, p.Class, p.Style, p.Title, p.Attributes)
	return elementNode{name: "script", props: a, void: false, children: children}
}
