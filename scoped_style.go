package html

import "strings"

// ScopedStyle renders CSS in a style tag, prefixing selectors in ordinary
// rules with the supplied id. At-rules are preserved and their contents are
// scanned recursively so nested rules are scoped too.
func ScopedStyle(scopeID string, css string) Node {
	return rawNode("<style>" + scopeCSS(scopeID, css, false) + "</style>")
}

func scopeCSS(scopeID, css string, inAtRule bool) string {
	var b strings.Builder
	for len(css) > 0 {
		i := strings.IndexByte(css, '{')
		if i < 0 {
			b.WriteString(css)
			break
		}
		header := css[:i]
		trimmed := strings.TrimSpace(header)
		at := inAtRule || strings.HasPrefix(trimmed, "@")
		if at {
			b.WriteString(header)
		} else {
			parts := strings.Split(header, ",")
			for j, p := range parts {
				if j > 0 {
					b.WriteByte(',')
				}
				b.WriteString(" #" + scopeID + " " + strings.TrimSpace(p))
			}
		}
		b.WriteByte('{')
		rest := css[i+1:]
		depth, close := 1, -1
		for j, r := range rest {
			if r == '{' {
				depth++
			} else if r == '}' {
				depth--
				if depth == 0 {
					close = j
					break
				}
			}
		}
		if close < 0 {
			b.WriteString(scopeCSS(scopeID, rest, at))
			break
		}
		b.WriteString(scopeCSS(scopeID, rest[:close], at))
		b.WriteByte('}')
		css = rest[close+1:]
	}
	return b.String()
}
