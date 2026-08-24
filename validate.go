package html

type validation struct {
	warnings []Warning
	ids      map[string]bool
}

func Validate(n Node) []Warning {
	v := &validation{ids: map[string]bool{}}
	n.validateNode(v, "")
	return v.warnings
}
func (v *validation) warn(code, msg, tag, path string) {
	v.warnings = append(v.warnings, Warning{Code: code, Message: msg, Tag: tag, Path: path})
}
