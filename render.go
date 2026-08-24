package html

import (
	"html"
	"strconv"
)

func Render(n Node) string {
	for _, w := range Validate(n) {
		emitWarning(w)
	}
	r := &renderer{}
	n.renderNode(r)
	return r.b.String()
}
func (r *renderer) text(s string) { r.b.WriteString(html.EscapeString(s)) }
func (r *renderer) raw(s string)  { r.b.WriteString(s) }
func (r *renderer) element(n elementNode) {
	r.b.WriteByte('<')
	r.b.WriteString(n.name)
	for _, a := range n.props {
		if a.Boolean {
			if a.Value == "true" {
				r.b.WriteByte(' ')
				r.b.WriteString(a.Name)
			}
			continue
		}
		r.b.WriteByte(' ')
		r.b.WriteString(a.Name)
		r.b.WriteString(`="`)
		r.b.WriteString(html.EscapeString(a.Value))
		r.b.WriteByte('"')
	}
	r.b.WriteByte('>')
	if n.void {
		return
	}
	for _, c := range n.children {
		c.renderNode(r)
	}
	r.b.WriteString(`</` + n.name + `>`)
}
func itoa(i int) string { return strconv.Itoa(i) }
