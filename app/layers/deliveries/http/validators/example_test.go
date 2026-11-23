package validators_test

import (
	"fmt"

	"github.com/chaihanij/backend-challenge/app/layers/deliveries/http/validators"
	"github.com/go-playground/validator/v10"
)

func ExampleMapValidationErrorsSlice() {
	type Request struct {
		Email string `json:"email" validate:"required,email"`
	}

	v := validator.New()
	validators.RegisterJSONTagName(v)

	err := v.Struct(Request{Email: ""})
	items := validators.MapValidationErrorsSlice(err)

	for _, item := range items {
		fmt.Printf("field=%s tag=%s\n", item.Field, item.Tag)
	}
	// Output:
	// field=email tag=required
}

func ExampleMapValidationErrorsSlice_nilError() {
	// nil error returns an empty (non-nil) slice.
	items := validators.MapValidationErrorsSlice(nil)
	fmt.Println(len(items))
	// Output:
	// 0
}
