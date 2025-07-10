package main

import "testing"

func TestFindFirstAvailableAncestor(t *testing.T) {
	tests := []struct {
		name               string
		parentID           *string
		unavailableParents map[string]*string
		expected           *string
	}{
		{
			name:               "nil parent",
			parentID:           nil,
			unavailableParents: map[string]*string{},
			expected:           nil,
		},
		{
			name:     "available parent",
			parentID: strPtr("parent1"),
			unavailableParents: map[string]*string{
				"parent2": strPtr("grandparent"),
			},
			expected: strPtr("parent1"),
		},
		{
			name:     "skip one unavailable",
			parentID: strPtr("parent1"),
			unavailableParents: map[string]*string{
				"parent1": strPtr("grandparent1"),
			},
			expected: strPtr("grandparent1"),
		},
		{
			name:     "skip multiple unavailable",
			parentID: strPtr("parent1"),
			unavailableParents: map[string]*string{
				"parent1":      strPtr("grandparent1"),
				"grandparent1": strPtr("great-grandparent1"),
			},
			expected: strPtr("great-grandparent1"),
		},
		{
			name:     "all ancestors unavailable",
			parentID: strPtr("parent1"),
			unavailableParents: map[string]*string{
				"parent1":      strPtr("grandparent1"),
				"grandparent1": nil,
			},
			expected: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := findFirstAvailableAncestor(tt.parentID, tt.unavailableParents)
			if !equalStringPtr(result, tt.expected) {
				t.Errorf("expected %v, got %v", strPtrToStr(tt.expected), strPtrToStr(result))
			}
		})
	}
}

// Helper functions for tests
func strPtr(s string) *string {
	return &s
}

func equalStringPtr(a, b *string) bool {
	if a == nil && b == nil {
		return true
	}
	if a == nil || b == nil {
		return false
	}
	return *a == *b
}

func strPtrToStr(p *string) string {
	if p == nil {
		return "<nil>"
	}
	return *p
}
