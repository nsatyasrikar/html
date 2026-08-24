package html

func UnsafeHTML(s string) Node { return rawNode(s) }
