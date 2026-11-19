package errors

import (
	"fmt"
)

func firstErr(err []error) error {
	if len(err) > 0 {
		return err[0]
	}
	return nil
}

func problemType(path string) string {
	return fmt.Sprintf("%s/problems/%s", BaseUrl, path)
}
