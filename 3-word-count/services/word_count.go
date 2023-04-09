package services

import (
	"strings"
	"word-count/models"
)

func WordCount(text string) models.WordCount {

	// Remove commas and periods
	text = strings.ReplaceAll(text, ",", "")
	text = strings.ReplaceAll(text, ".", "")

	// Split the text into words
	words := strings.Fields(text)

	// Count the words
	wordCount := make(map[string]int)
	for _, word := range words {
		wordCount[word]++
	}

	// Place the word count inside struct
	output := models.WordCount{
		Beef: wordCount,
	}

	return output
}
