package handlers

import (
	"encoding/json"
	"fmt"
	"io/ioutil"

	"github.com/labstack/gommon/log"
)

func readFile(path string) ([][]int, error) {
	// Define a 2D slice to hold the matrix
	var matrix [][]int

	// Read the JSON file into a byte slice
	bytes, err := ioutil.ReadFile(path)
	if err != nil {
		log.Error(fmt.Sprintf("Cannot read file got error: %v", err))
		return matrix, err
	}

	// Unmarshal the JSON data into the matrix slice
	err = json.Unmarshal(bytes, &matrix)
	if err != nil {
		log.Error(fmt.Sprintf("Cannot unmarshal matrix got error: %v", err))
		return matrix, err
	}
	return matrix, nil
}
