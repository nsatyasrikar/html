package html

type AudioProps struct {
	ID, Class, Style, Title, Src, Preload *string
	Controls                              *bool
	Autoplay                              *bool
	Loop                                  *bool
	Muted                                 *bool
	Attributes                            []Attribute
}

func Audio(p AudioProps, children ...Node) Node {
	a := globalAttrs(p.ID, p.Class, p.Style, p.Title, p.Attributes)
	if p.Src != nil {
		a = append(a, Attribute{"src", *p.Src, false})
	}
	if p.Preload != nil {
		a = append(a, Attribute{"preload", *p.Preload, false})
	}
	if p.Controls != nil {
		a = append(a, Attribute{"controls", "true", true})
	}
	if p.Autoplay != nil {
		a = append(a, Attribute{"autoplay", "true", true})
	}
	if p.Loop != nil {
		a = append(a, Attribute{"loop", "true", true})
	}
	if p.Muted != nil {
		a = append(a, Attribute{"muted", "true", true})
	}
	return elementNode{name: "audio", props: a, void: false, children: children}
}
