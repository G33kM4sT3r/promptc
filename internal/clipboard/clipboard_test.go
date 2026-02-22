package clipboard

import "testing"

func TestAvailable(t *testing.T) {
	// Should not panic regardless of platform
	_ = Available()
}

func TestCopy_EmptyString(t *testing.T) {
	if !Available() {
		t.Skip("clipboard not available on this platform")
	}
	err := Copy("")
	if err != nil {
		t.Errorf("Copy empty string: %v", err)
	}
}

func TestCopy_NonEmpty(t *testing.T) {
	if !Available() {
		t.Skip("clipboard not available on this platform")
	}
	err := Copy("hello promptc")
	if err != nil {
		t.Errorf("Copy: %v", err)
	}
}
