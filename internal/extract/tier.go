package extract

import "promptc/internal/slots"

// CalculateTier returns the semantic depth tier based on extracted slot values.
// Returns "minimal", "standard", or "rich". Deterministic: same slots = same tier.
func CalculateTier(s slots.Slots) string {
	if s.Depth == "short" {
		return "minimal"
	}
	if s.Depth == "deep" ||
		s.Audience == "advanced" ||
		len(s.Entities) >= 2 ||
		s.Stage == "optimization" {
		return "rich"
	}
	return "standard"
}
