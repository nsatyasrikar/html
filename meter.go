package html

type MeterProps struct {
	ID, Class, Style, Title, Value, Min, Max, Low, High, Optimum *string
	Attributes                                                   []Attribute
}

func Meter(p MeterProps, children ...Node) Node {
	a := globalAttrs(p.ID, p.Class, p.Style, p.Title, p.Attributes)
	if p.Value != nil {
		a = append(a, Attribute{"value", *p.Value, false})
	}
	if p.Min != nil {
		a = append(a, Attribute{"min", *p.Min, false})
	}
	if p.Max != nil {
		a = append(a, Attribute{"max", *p.Max, false})
	}
	if p.Low != nil {
		a = append(a, Attribute{"low", *p.Low, false})
	}
	if p.High != nil {
		a = append(a, Attribute{"high", *p.High, false})
	}
	if p.Optimum != nil {
		a = append(a, Attribute{"optimum", *p.Optimum, false})
	}
	return elementNode{name: "meter", props: a, void: false, children: children}
}
