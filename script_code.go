package html

import "strings"

type scriptCodeNode string

func ScriptCode(source string) Node                         { return scriptCodeNode(source) }
func (n scriptCodeNode) renderNode(r *renderer)             { r.raw(escapeScriptEnd(string(n))) }
func (scriptCodeNode) validateNode(_ *validation, _ string) {}

func escapeScriptEnd(s string) string {
	if !strings.Contains(strings.ToLower(s), "</script") {
		return s
	}
	b := strings.Builder{}
	b.Grow(len(s) + 8)
	for i := 0; i < len(s); i++ {
		if s[i] == '<' && i+7 < len(s) && s[i+1] == '/' && strings.EqualFold(s[i+2:i+8], "script") {
			b.WriteString("<\\/")
			i++
			continue
		}
		b.WriteByte(s[i])
	}
	return b.String()
}
