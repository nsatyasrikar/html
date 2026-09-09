package html

type MarqueeProps struct {
	ID, Class, Style, Title *string
	Attributes              []Attribute
}

func Marquee(p MarqueeProps, children ...Node) Node {
	a := globalAttrs(p.ID, p.Class, p.Style, p.Title, p.Attributes)
	return elementNode{name: "marquee", props: a, void: false, children: children}
}
