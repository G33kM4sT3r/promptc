package tokenize

import (
	"reflect"
	"testing"
)

func TestTokenize(t *testing.T) {
	tests := []struct {
		input string
		want  []string
	}{
		{"Hello World", []string{"hello", "world"}},
		{"EXPLAIN Closures", []string{"explain", "closures"}},
		{"  spaced  out  ", []string{"spaced", "out"}},
		{"Node.js API", []string{"node.js", "api"}},
		{"", []string{}},
	}

	for _, tt := range tests {
		got := Tokenize(tt.input)
		if !reflect.DeepEqual(got, tt.want) {
			t.Errorf("Tokenize(%q) = %v, want %v", tt.input, got, tt.want)
		}
	}
}
