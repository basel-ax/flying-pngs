package ui

import (
	"testing"
)

// TestParseInt tests the parseInt helper function
func TestParseInt(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		min      int
		max      int
		def      int
		expected int
	}{
		{
			name:     "valid number",
			input:    "100",
			min:      0,
			max:      200,
			def:      50,
			expected: 100,
		},
		{
			name:     "below minimum",
			input:    "10",
			min:      50,
			max:      200,
			def:      100,
			expected: 50,
		},
		{
			name:     "above maximum",
			input:    "300",
			min:      0,
			max:      200,
			def:      100,
			expected: 200,
		},
		{
			name:     "invalid number",
			input:    "abc",
			min:      0,
			max:      200,
			def:      50,
			expected: 50,
		},
		{
			name:     "empty string",
			input:    "",
			min:      0,
			max:      200,
			def:      50,
			expected: 50,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := parseInt(tt.input, tt.min, tt.max, tt.def); got != tt.expected {
				t.Errorf("parseInt() = %v, want %v", got, tt.expected)
			}
		})
	}
}
