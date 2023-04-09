package services

import "testing"

func Test_maxSumPath(t *testing.T) {
	testCases := []struct {
		name     string
		matrix   [][]int
		expected int
	}{
		{
			name: "General",
			matrix: [][]int{
				{59},
				{73, 41},
				{52, 40, 53},
				{26, 53, 6, 34},
			},
			expected: 237,
		},
		{
			name: "Single element",
			matrix: [][]int{
				{5},
			},
			expected: 5,
		},
		{
			name: "All positive elements",
			matrix: [][]int{
				{1},
				{2, 3},
				{4, 5, 6},
			},
			expected: 10,
		},
		{
			name: "All negative elements",
			matrix: [][]int{
				{-1},
				{-2, -3},
				{-4, -5, -6},
			},
			expected: -7,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := MaxSumPath(tc.matrix)
			if result != tc.expected {
				t.Errorf("Expected %d, got %d", tc.expected, result)
			}
		})
	}
}
