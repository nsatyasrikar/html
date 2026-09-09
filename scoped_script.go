package html

// ScopedScript wraps JavaScript in an IIFE to keep declarations local.
func ScopedScript(js string) Node {
	return rawNode("<script>(function(){\n" + js + "\n})();</script>")
}
