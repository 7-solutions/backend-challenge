// Package validators bridges github.com/go-playground/validator/v10 with the
// library's error and response types.
//
// # Startup registration
//
// Call RegisterJSONTagName once at application startup on your validator
// instance. This configures the validator to report field names using the json
// struct tag (e.g. "email") rather than the Go field name (e.g. "Email"):
//
//	v := validator.New()
//	validators.RegisterJSONTagName(v)
//
// # Parsing validation errors
//
// MapValidationErrorsSlice converts a validator.ValidationErrors error into a
// slice of FieldErrorItem, each with Field (dot-path), Message (human-readable),
// Tag (rule name), and Value (failing value):
//
//	items := validators.MapValidationErrorsSlice(err)
//
// The function always returns a non-nil slice, so it is safe to JSON-marshal
// even when err is nil or is not a validation error (returns an empty slice).
//
// # Built-in message templates
//
// Human-readable messages are generated for: required, email, url, uuid,
// alphanum, min, max, len, gte, lte, oneof, eqfield. All other tags fall back
// to "is invalid".
//
// # Integration with ProblemDetails
//
// The typical error-handling flow in a Gin handler:
//
//	if err := c.ShouldBindJSON(&req); err != nil {
//	    items := validators.MapValidationErrorsSlice(err)
//	    fields := make(map[string]string, len(items))
//	    for _, item := range items {
//	        fields[item.Field] = item.Message
//	    }
//	    _ = c.Error(httpxerrors.NewBadRequest("validation failed",
//	        cmerrors.NewValidationError(fields)))
//	    return
//	}
//
// Alternatively, pass the raw error directly to NewBadRequest — it applies
// MapValidationErrorsSlice internally.
package validators
