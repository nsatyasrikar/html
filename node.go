package html

import "strings"

type Node interface {
	renderNode(*renderer)
	validateNode(*validation, string)
}

type textNode string

func (n textNode) renderNode(r *renderer)               { r.text(string(n)) }
func (n textNode) validateNode(_ *validation, _ string) {}

type rawNode string

func (n rawNode) renderNode(r *renderer) { r.raw(string(n)) }
func (rawNode) validateNode(v *validation, path string) {
	v.warn("unsafe-html", "unsafe HTML emitted", "", path)
}

type elementNode struct {
	name       string
	props      []Attribute
	children   []Node
	void       bool
	deprecated bool
}

func (n elementNode) renderNode(r *renderer) { r.element(n) }
func (n elementNode) validateNode(v *validation, path string) {
	if n.void && len(n.children) > 0 {
		v.warn("void-children", "void element has children", n.name, path)
	}
	seen := map[string]bool{}
	for _, a := range n.props {
		k := a.Name
		if seen[k] {
			v.warn("duplicate-attribute", "duplicate attribute", n.name, path)
		}
		seen[k] = true
		if k == "id" && a.Value != "" {
			if v.ids[a.Value] {
				v.warn("duplicate-id", "duplicate id", n.name, path)
			}
			v.ids[a.Value] = true
		}
	}
	if n.deprecated {
		v.warn("deprecated", "deprecated element", n.name, path)
	}
	for i, c := range n.children {
		c.validateNode(v, path+"/"+n.name+"["+itoa(i)+"]")
	}
}

type renderer struct{ b strings.Builder }
