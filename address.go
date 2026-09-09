package html

type AddressProps struct {
	ID, Class, Style, Title *string
	Attributes              []Attribute
}

func Address(p AddressProps, children ...Node) Node {
	a := globalAttrs(p.ID, p.Class, p.Style, p.Title, p.Attributes)
	return elementNode{name: "address", props: a, void: false, children: children}
}
