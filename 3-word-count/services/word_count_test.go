package services

import (
	"reflect"
	"testing"
	"word-count/models"
)

func Test_WordCount(t *testing.T) {
	text := "Consequat venison quis jowl porchetta minim ea flank, burgdoggen adipisicing. Non ex in sint doner, jowl boudin quis pariatur tongue tri-tip proident."
	testCases := []struct {
		name     string
		input    string
		expected models.WordCount
	}{
		{
			name:  "General case",
			input: text,
			expected: models.WordCount{
				Beef: map[string]int{
					"Consequat": 1, "venison": 1, "quis": 2, "jowl": 2, "porchetta": 1, "minim": 1, "ea": 1, "flank": 1, "burgdoggen": 1, "adipisicing": 1,
					"Non": 1, "ex": 1, "in": 1, "sint": 1, "doner": 1, "boudin": 1, "pariatur": 1, "tongue": 1, "tri-tip": 1, "proident": 1,
				},
			},
		},
		{
			name:  "One word",
			input: "เนื้อ",
			expected: models.WordCount{
				Beef: map[string]int{
					"เนื้อ": 1,
				},
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := WordCount(tc.input)
			if !reflect.DeepEqual(result, tc.expected) {
				t.Errorf("Expected %v, but got %v", tc.expected, result)
			}
		})
	}
}
