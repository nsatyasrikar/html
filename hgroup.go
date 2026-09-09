package html

type HgroupProps struct {
	ID, Class, Style, Title *string
	Attributes              []Attribute
}

func Hgroup(p HgroupProps, children ...Node) Node {
	a := globalAttrs(p.ID, p.Class, p.Style, p.Title, p.Attributes)
	return elementNode{name: "hgroup", props: a, void: false, children: children}
}
