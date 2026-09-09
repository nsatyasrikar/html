package html

type TrackProps struct {
	ID, Class, Style, Title *string
	Attributes              []Attribute
}

func Track(p TrackProps) Node {
	a := globalAttrs(p.ID, p.Class, p.Style, p.Title, p.Attributes)
	return elementNode{name: "track", props: a, void: true}
}
