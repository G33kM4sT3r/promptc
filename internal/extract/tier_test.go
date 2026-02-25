package extract

import (
	"testing"

	"promptc/internal/slots"
)

func TestCalculateTier(t *testing.T) {
	tests := []struct {
		name string
		s    slots.Slots
		want string
	}{
		{"default_standard", slots.Slots{Intent: "explain"}, "standard"},
		{"short_forces_minimal", slots.Slots{Depth: "short"}, "minimal"},
		{"short_overrides_advanced", slots.Slots{Depth: "short", Audience: "advanced"}, "minimal"},
		{"short_overrides_multi_entity", slots.Slots{
			Depth:    "short",
			Entities: []slots.Entity{{Text: "a", Role: "x"}, {Text: "b", Role: "y"}},
		}, "minimal"},
		{"deep_triggers_rich", slots.Slots{Depth: "deep"}, "rich"},
		{"advanced_triggers_rich", slots.Slots{Audience: "advanced"}, "rich"},
		{"two_entities_triggers_rich", slots.Slots{
			Entities: []slots.Entity{{Text: "a", Role: "x"}, {Text: "b", Role: "y"}},
		}, "rich"},
		{"optimization_triggers_rich", slots.Slots{Stage: "optimization"}, "rich"},
		{"beginner_standard", slots.Slots{Audience: "beginners"}, "standard"},
		{"one_entity_standard", slots.Slots{
			Entities: []slots.Entity{{Text: "a", Role: "x"}},
		}, "standard"},
		{"empty_slots_standard", slots.Slots{}, "standard"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := CalculateTier(tt.s)
			if got != tt.want {
				t.Errorf("CalculateTier() = %q, want %q", got, tt.want)
			}
		})
	}
}
