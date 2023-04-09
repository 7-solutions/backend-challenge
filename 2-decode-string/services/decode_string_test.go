package services

import (
	"reflect"
	"testing"
)

func Test_DecodeV2(t *testing.T) {
	testCases := []struct {
		name     string
		input    string
		expected []int
	}{
		{
			name:     "Empty string",
			input:    "",
			expected: []int{},
		},
		{
			name:     "Single letter L",
			input:    "L",
			expected: []int{1, 0},
		},
		{
			name:     "Single letter R",
			input:    "R",
			expected: []int{0, 1},
		},
		{
			name:     "Two letter LR",
			input:    "LR",
			expected: []int{1, 0, 1},
		},
		{
			name:     "Two letter RL",
			input:    "RL",
			expected: []int{0, 1, 0},
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := DecodeV2(tc.input)
			if !reflect.DeepEqual(result, tc.expected) {
				t.Errorf("Expected %v, but got %v", tc.expected, result)
			}
		})
	}
}
