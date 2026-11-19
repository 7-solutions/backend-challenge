package validators

import (
	"errors"
	"fmt"
	"reflect"
	"strings"

	"github.com/go-playground/validator/v10"
)

// FieldErrorItem describes a single field-level validation failure.
type FieldErrorItem struct {
	// Field is the dot-separated JSON path to the failing field (e.g. "profile.email").
	Field string `json:"field"`
	// Message is a human-readable description of the validation failure.
	Message string `json:"message"`
	// Tag is the validator rule that failed (e.g. "required", "email").
	Tag string `json:"tag,omitempty"`
	// Value is the value that failed validation.
	Value any `json:"value,omitempty"`
}

// FieldErrorItems is a slice of FieldErrorItem with a helper for map conversion.
type FieldErrorItems []FieldErrorItem

// GetFieldErrorItems converts the slice to a map of field → message suitable
// for passing to cmerrors.NewValidationError.
func (f FieldErrorItems) GetFieldErrorItems() map[string]string {
	out := make(map[string]string, len(f))
	for _, item := range f {
		out[item.Field] = item.Message
	}
	return out
}

// RegisterJSONTagName configures the validator to report field names using
// the `json` struct tag rather than Go field names. Call this once at startup.
func RegisterJSONTagName(v *validator.Validate) {
	v.RegisterTagNameFunc(func(fld reflect.StructField) string {
		name := strings.SplitN(fld.Tag.Get("json"), ",", 2)[0]
		if name == "-" {
			return ""
		}
		return name
	})
}

// MapValidationErrorsSlice converts a validator.ValidationErrors error to a
// slice of FieldErrorItem. It always returns a non-nil slice, so the result is
// safe to JSON-marshal even when err is nil or is not a validation error
// (returns an empty slice in those cases).
func MapValidationErrorsSlice(err error) []FieldErrorItem {
	out := make([]FieldErrorItem, 0) // non-nil so JSON marshals as []

	if err == nil {
		return out
	}

	var verrs validator.ValidationErrors
	if !errors.As(err, &verrs) {
		return out
	}

	for _, fe := range verrs {
		out = append(out, FieldErrorItem{
			Field:   fieldPath(fe),
			Message: validationMessage(fe),
			Tag:     fe.Tag(),
			Value:   fe.Value(),
		})
	}
	return out
}

// fieldPath returns the dot-separated path to the failing field, excluding
// the root struct name. e.g. "profile.email" not "CreateUserReq.profile.email".
func fieldPath(fe validator.FieldError) string {
	ns := fe.Namespace()
	if idx := strings.Index(ns, "."); idx != -1 {
		return ns[idx+1:]
	}
	return fe.Field()
}

func validationMessage(fe validator.FieldError) string {
	if b, ok := messageBuilders[fe.Tag()]; ok {
		return b(fe)
	}
	return "is invalid"
}

type messageBuilder func(validator.FieldError) string

var messageBuilders = map[string]messageBuilder{
	"required": staticMsg("is required"),
	"email":    staticMsg("must be a valid email"),
	"url":      staticMsg("must be a valid URL"),
	"uuid":     staticMsg("must be a valid UUID"),
	"alphanum": staticMsg("must be alphanumeric"),
	"min":      minMessage,
	"max":      maxMessage,
	"len":      lenMessage,
	"gte":      func(fe validator.FieldError) string { return fmt.Sprintf("must be ≥ %s", fe.Param()) },
	"lte":      func(fe validator.FieldError) string { return fmt.Sprintf("must be ≤ %s", fe.Param()) },
	"oneof": func(fe validator.FieldError) string {
		opts := strings.ReplaceAll(fe.Param(), " ", ", ")
		return fmt.Sprintf("must be one of: %s", opts)
	},
	"eqfield": func(fe validator.FieldError) string {
		return fmt.Sprintf("must match %s", fe.Param())
	},
}

func staticMsg(msg string) messageBuilder {
	return func(validator.FieldError) string { return msg }
}

func minMessage(fe validator.FieldError) string {
	switch fe.Kind() {
	case reflect.String:
		return fmt.Sprintf("must be at least %s characters", fe.Param())
	case reflect.Slice, reflect.Array, reflect.Map:
		return fmt.Sprintf("must contain at least %s items", fe.Param())
	default:
		return fmt.Sprintf("must be at least %s", fe.Param())
	}
}

func maxMessage(fe validator.FieldError) string {
	switch fe.Kind() {
	case reflect.String:
		return fmt.Sprintf("must be at most %s characters", fe.Param())
	case reflect.Slice, reflect.Array, reflect.Map:
		return fmt.Sprintf("must contain at most %s items", fe.Param())
	default:
		return fmt.Sprintf("must be at most %s", fe.Param())
	}
}

func lenMessage(fe validator.FieldError) string {
	switch fe.Kind() {
	case reflect.String:
		return fmt.Sprintf("must be exactly %s characters", fe.Param())
	case reflect.Slice, reflect.Array, reflect.Map:
		return fmt.Sprintf("must contain exactly %s items", fe.Param())
	default:
		return fmt.Sprintf("must equal %s", fe.Param())
	}
}
