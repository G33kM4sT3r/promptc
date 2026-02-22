package slots

import "testing"

func TestSlots_ZeroValue(t *testing.T) {
	var s Slots
	if s.Language != "" || s.Intent != "" || s.Topic != "" {
		t.Error("zero value should have empty strings")
	}
	if s.Entities != nil {
		t.Error("zero value should have nil entities")
	}
}

func TestEntity_Fields(t *testing.T) {
	e := Entity{Text: "Python", Role: "implementation_medium"}
	if e.Text != "Python" {
		t.Errorf("expected Python, got %s", e.Text)
	}
	if e.Role != "implementation_medium" {
		t.Errorf("expected implementation_medium, got %s", e.Role)
	}
}
