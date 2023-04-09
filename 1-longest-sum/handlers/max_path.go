package handlers

import (
	"longest-sum/services"

	"github.com/labstack/gommon/log"
)

func SumMaxPath(path string) (int, error) {

	var result int

	matrix, err := readFile(path)
	if err != nil {
		return result, err
	}

	log.Info("Successfully read file: %v", matrix)

	result = services.MaxSumPath(matrix)

	return result, nil
}
