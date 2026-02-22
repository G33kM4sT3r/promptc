package prompt

import "slices"

// AppendUnique appends items to slice only if not already present.
func AppendUnique(slice []string, items ...string) []string {
	for _, item := range items {
		found := slices.Contains(slice, item)
		if !found {
			slice = append(slice, item)
		}
	}
	return slice
}
