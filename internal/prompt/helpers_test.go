package prompt

import "testing"

func TestAppendUnique(t *testing.T) {
	slice := []string{"a", "b"}
	slice = AppendUnique(slice, "c")
	if len(slice) != 3 {
		t.Errorf("expected 3, got %d", len(slice))
	}
	slice = AppendUnique(slice, "b") // duplicate
	if len(slice) != 3 {
		t.Errorf("expected 3 (no dup), got %d", len(slice))
	}
	slice = AppendUnique(slice, "d", "a") // one new, one dup
	if len(slice) != 4 {
		t.Errorf("expected 4, got %d", len(slice))
	}
}

func TestAppendUniqueEmpty(t *testing.T) {
	var slice []string
	slice = AppendUnique(slice, "a")
	if len(slice) != 1 {
		t.Errorf("expected 1, got %d", len(slice))
	}
}
