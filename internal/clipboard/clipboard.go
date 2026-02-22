package clipboard

import "github.com/atotto/clipboard"

// Copy writes text to the system clipboard.
func Copy(text string) error {
	return clipboard.WriteAll(text)
}

// Available reports whether clipboard access is likely to work.
func Available() bool {
	return !clipboard.Unsupported
}
