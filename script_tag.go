package html

type ScriptTagProps struct {
	ID, Class, Style, Title, Src, Type *string
	Async                              *bool
	Defer                              *bool
	Attributes                         []Attribute
}

func ScriptTag(p ScriptTagProps, children ...Node) Node {
	a := globalAttrs(p.ID, p.Class, p.Style, p.Title, p.Attributes)
	if p.Src != nil {
		a = append(a, Attribute{"src", *p.Src, false})
	}
	if p.Type != nil {
		a = append(a, Attribute{"type", *p.Type, false})
	}
	if p.Async != nil {
		a = append(a, Attribute{"async", "true", true})
	}
	if p.Defer != nil {
		a = append(a, Attribute{"defer", "true", true})
	}
	return elementNode{name: "script", props: a, void: false, children: children}
}
