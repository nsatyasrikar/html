package html

// Fragment groups nodes so a function can build a []Node in pieces
// and pass it as a single children argument, without an extra wrapper element.
func Fragment(children ...Node) []Node {
	return children
}
