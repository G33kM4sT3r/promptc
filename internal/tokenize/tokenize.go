package tokenize

import "strings"

func Tokenize(input string) []string {
	return strings.Fields(strings.ToLower(input))
}
