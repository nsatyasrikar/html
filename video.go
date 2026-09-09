package html

type VideoProps struct {
	ID, Class, Style, Title, Src, Poster, Width, Height, Preload *string
	Controls                                                     *bool
	Autoplay                                                     *bool
	Loop                                                         *bool
	Muted                                                        *bool
	Attributes                                                   []Attribute
}

func Video(p VideoProps, children ...Node) Node {
	a := globalAttrs(p.ID, p.Class, p.Style, p.Title, p.Attributes)
	if p.Src != nil {
		a = append(a, Attribute{"src", *p.Src, false})
	}
	if p.Poster != nil {
		a = append(a, Attribute{"poster", *p.Poster, false})
	}
	if p.Width != nil {
		a = append(a, Attribute{"width", *p.Width, false})
	}
	if p.Height != nil {
		a = append(a, Attribute{"height", *p.Height, false})
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
	return elementNode{name: "video", props: a, void: false, children: children}
}
