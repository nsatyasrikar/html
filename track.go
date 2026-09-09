package html

type TrackProps struct {
	ID, Class, Style, Title, Src, Kind, Srclang, Label *string
	Default                                            *bool
	Attributes                                         []Attribute
}

func Track(p TrackProps) Node {
	a := globalAttrs(p.ID, p.Class, p.Style, p.Title, p.Attributes)
	if p.Src != nil {
		a = append(a, Attribute{"src", *p.Src, false})
	}
	if p.Kind != nil {
		a = append(a, Attribute{"kind", *p.Kind, false})
	}
	if p.Srclang != nil {
		a = append(a, Attribute{"srclang", *p.Srclang, false})
	}
	if p.Label != nil {
		a = append(a, Attribute{"label", *p.Label, false})
	}
	if p.Default != nil {
		a = append(a, Attribute{"default", "true", true})
	}
	return elementNode{name: "track", props: a, void: true}
}
