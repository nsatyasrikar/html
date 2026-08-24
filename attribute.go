package html

import "regexp"

type Attribute struct {
	Name    string
	Value   string
	Boolean bool
}

var attrName = regexp.MustCompile(`^[A-Za-z_:][A-Za-z0-9_.:-]*$`)

func customAttributes(in []Attribute) []Attribute {
	out := make([]Attribute, 0, len(in))
	for _, a := range in {
		if attrName.MatchString(a.Name) {
			out = append(out, a)
		}
	}
	return out
}
