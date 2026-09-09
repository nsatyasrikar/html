package html

type MarqueeProps struct {
	ID, Class, Style, Title, Behavior, Direction, Scrollamount, Scrolldelay, Loop *string
	Attributes                                                                    []Attribute
}

func Marquee(p MarqueeProps, children ...Node) Node {
	a := globalAttrs(p.ID, p.Class, p.Style, p.Title, p.Attributes)
	if p.Behavior != nil {
		a = append(a, Attribute{"behavior", *p.Behavior, false})
	}
	if p.Direction != nil {
		a = append(a, Attribute{"direction", *p.Direction, false})
	}
	if p.Scrollamount != nil {
		a = append(a, Attribute{"scrollamount", *p.Scrollamount, false})
	}
	if p.Scrolldelay != nil {
		a = append(a, Attribute{"scrolldelay", *p.Scrolldelay, false})
	}
	if p.Loop != nil {
		a = append(a, Attribute{"loop", *p.Loop, false})
	}
	return elementNode{name: "marquee", props: a, void: false, children: children}
}
