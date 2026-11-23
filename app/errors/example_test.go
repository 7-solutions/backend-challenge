package errors_test

import (
	"errors"
	"fmt"

	cmerrors "github.com/chaihanij/backend-challenge/app/errors"
)

func ExampleNewNotFoundError() {
	err := cmerrors.NewNotFoundError("User", "42")
	fmt.Println(err)
	// Output:
	// User with ID 42 not found
}

func ExampleWrapNotFound() {
	cause := errors.New("sql: no rows")
	err := cmerrors.WrapNotFound("User", "42", cause)

	var nfe *cmerrors.NotFoundError
	if errors.As(err, &nfe) {
		fmt.Printf("resource=%s id=%s\n", nfe.Resource, nfe.ID)
	}
	// Output:
	// resource=User id=42
}

func ExampleNewValidationError() {
	err := cmerrors.NewValidationError(map[string]string{
		"email": "must be a valid email",
	})
	fmt.Println(err)
	// Output:
	// validation failed
}

func ExampleWrap() {
	base := errors.New("connection refused")
	wrapped := cmerrors.Wrap(base, "GetUser")
	fmt.Println(wrapped)

	// Wrap(nil, ...) returns nil.
	fmt.Println(cmerrors.Wrap(nil, "op"))
	// Output:
	// GetUser: connection refused
	// <nil>
}
