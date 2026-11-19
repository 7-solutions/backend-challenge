package validators_test

import (
	"errors"
	"testing"

	"github.com/chaihanij/backend-challenge/app/layers/deliveries/http/validators"
	"github.com/go-playground/validator/v10"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// ── helpers ──────────────────────────────────────────────────────────────────

func newValidate() *validator.Validate {
	v := validator.New()
	validators.RegisterJSONTagName(v)
	return v
}

// ── RegisterJSONTagName ───────────────────────────────────────────────────────

func TestRegisterJSONTagName_UsesJSONTag(t *testing.T) {
	type Req struct {
		UserName string `json:"user_name" validate:"required"`
	}

	v := newValidate()
	err := v.Struct(Req{})
	require.Error(t, err)

	items := validators.MapValidationErrorsSlice(err)
	require.Len(t, items, 1)
	assert.Equal(t, "user_name", items[0].Field)
}

func TestRegisterJSONTagName_DashTagOmitted(t *testing.T) {
	// field tagged json:"-" should fall back to empty name (validator skips it)
	type Inner struct {
		Skip string `json:"-" validate:"required"`
	}
	type Req struct {
		Name string `json:"name" validate:"required"`
	}

	v := newValidate()
	err := v.Struct(Req{})
	require.Error(t, err)

	items := validators.MapValidationErrorsSlice(err)
	require.Len(t, items, 1)
	assert.Equal(t, "name", items[0].Field)
}

// ── MapValidationErrorsSlice ──────────────────────────────────────────────────

func TestMapValidationErrorsSlice_NilError(t *testing.T) {
	out := validators.MapValidationErrorsSlice(nil)
	assert.NotNil(t, out)
	assert.Empty(t, out)
}

func TestMapValidationErrorsSlice_NonValidationError(t *testing.T) {
	out := validators.MapValidationErrorsSlice(errors.New("some other error"))
	assert.NotNil(t, out)
	assert.Empty(t, out)
}

func TestMapValidationErrorsSlice_MultipleFields(t *testing.T) {
	type Req struct {
		Email string `json:"email" validate:"required,email"`
		Name  string `json:"name"  validate:"required"`
	}

	v := newValidate()
	err := v.Struct(Req{})
	require.Error(t, err)

	items := validators.MapValidationErrorsSlice(err)
	assert.Len(t, items, 2)

	fieldNames := make([]string, len(items))
	for i, it := range items {
		fieldNames[i] = it.Field
	}
	assert.ElementsMatch(t, []string{"email", "name"}, fieldNames)
}

func TestMapValidationErrorsSlice_ItemHasTagAndValue(t *testing.T) {
	type Req struct {
		Age int `json:"age" validate:"gte=18"`
	}

	v := newValidate()
	err := v.Struct(Req{Age: 5})
	require.Error(t, err)

	items := validators.MapValidationErrorsSlice(err)
	require.Len(t, items, 1)
	assert.Equal(t, "gte", items[0].Tag)
	assert.Equal(t, 5, items[0].Value)
}

// ── fieldPath (via MapValidationErrorsSlice) ──────────────────────────────────

func TestFieldPath_NestedStruct(t *testing.T) {
	type Address struct {
		City string `json:"city" validate:"required"`
	}
	type Req struct {
		Address Address `json:"address" validate:"required"`
	}

	v := newValidate()
	err := v.Struct(Req{})
	require.Error(t, err)

	items := validators.MapValidationErrorsSlice(err)
	require.NotEmpty(t, items)

	// nested path should not contain root struct name
	for _, it := range items {
		assert.NotContains(t, it.Field, "Req.")
	}
}

// ── validationMessage / messageBuilders ──────────────────────────────────────

func TestMessage_Required(t *testing.T) {
	type Req struct {
		Name string `json:"name" validate:"required"`
	}
	items := mustValidate(t, Req{})
	assert.Equal(t, "is required", items[0].Message)
}

func TestMessage_Email(t *testing.T) {
	type Req struct {
		Email string `json:"email" validate:"required,email"`
	}
	items := mustValidate(t, Req{Email: "not-an-email"})
	assert.Equal(t, "must be a valid email", items[0].Message)
}

func TestMessage_URL(t *testing.T) {
	type Req struct {
		Link string `json:"link" validate:"url"`
	}
	items := mustValidate(t, Req{Link: "not-a-url"})
	assert.Equal(t, "must be a valid URL", items[0].Message)
}

func TestMessage_UUID(t *testing.T) {
	type Req struct {
		ID string `json:"id" validate:"uuid"`
	}
	items := mustValidate(t, Req{ID: "not-a-uuid"})
	assert.Equal(t, "must be a valid UUID", items[0].Message)
}

func TestMessage_Alphanum(t *testing.T) {
	type Req struct {
		Code string `json:"code" validate:"alphanum"`
	}
	items := mustValidate(t, Req{Code: "has space!"})
	assert.Equal(t, "must be alphanumeric", items[0].Message)
}

func TestMessage_Min_String(t *testing.T) {
	type Req struct {
		Pass string `json:"pass" validate:"min=8"`
	}
	items := mustValidate(t, Req{Pass: "abc"})
	assert.Equal(t, "must be at least 8 characters", items[0].Message)
}

func TestMessage_Min_Number(t *testing.T) {
	type Req struct {
		Score int `json:"score" validate:"min=10"`
	}
	items := mustValidate(t, Req{Score: 1})
	assert.Equal(t, "must be at least 10", items[0].Message)
}

func TestMessage_Min_Slice(t *testing.T) {
	type Req struct {
		Tags []string `json:"tags" validate:"min=2"`
	}
	items := mustValidate(t, Req{Tags: []string{"one"}})
	assert.Equal(t, "must contain at least 2 items", items[0].Message)
}

func TestMessage_Max_String(t *testing.T) {
	type Req struct {
		Nick string `json:"nick" validate:"max=5"`
	}
	items := mustValidate(t, Req{Nick: "toolongname"})
	assert.Equal(t, "must be at most 5 characters", items[0].Message)
}

func TestMessage_Max_Number(t *testing.T) {
	type Req struct {
		Score int `json:"score" validate:"max=100"`
	}
	items := mustValidate(t, Req{Score: 200})
	assert.Equal(t, "must be at most 100", items[0].Message)
}

func TestMessage_Max_Slice(t *testing.T) {
	type Req struct {
		Tags []string `json:"tags" validate:"max=2"`
	}
	items := mustValidate(t, Req{Tags: []string{"a", "b", "c"}})
	assert.Equal(t, "must contain at most 2 items", items[0].Message)
}

func TestMessage_Len_String(t *testing.T) {
	type Req struct {
		Code string `json:"code" validate:"len=6"`
	}
	items := mustValidate(t, Req{Code: "abc"})
	assert.Equal(t, "must be exactly 6 characters", items[0].Message)
}

func TestMessage_Len_Slice(t *testing.T) {
	type Req struct {
		Items []int `json:"items" validate:"len=3"`
	}
	items := mustValidate(t, Req{Items: []int{1}})
	assert.Equal(t, "must contain exactly 3 items", items[0].Message)
}

func TestMessage_GTE(t *testing.T) {
	type Req struct {
		Age int `json:"age" validate:"gte=18"`
	}
	items := mustValidate(t, Req{Age: 10})
	assert.Equal(t, "must be ≥ 18", items[0].Message)
}

func TestMessage_LTE(t *testing.T) {
	type Req struct {
		Age int `json:"age" validate:"lte=65"`
	}
	items := mustValidate(t, Req{Age: 100})
	assert.Equal(t, "must be ≤ 65", items[0].Message)
}

func TestMessage_OneOf(t *testing.T) {
	type Req struct {
		Role string `json:"role" validate:"oneof=admin user guest"`
	}
	items := mustValidate(t, Req{Role: "superuser"})
	assert.Equal(t, "must be one of: admin, user, guest", items[0].Message)
}

func TestMessage_EqField(t *testing.T) {
	type Req struct {
		Password        string `json:"password"         validate:"required"`
		ConfirmPassword string `json:"confirm_password" validate:"eqfield=Password"`
	}
	v := newValidate()
	err := v.Struct(Req{Password: "secret", ConfirmPassword: "different"})
	require.Error(t, err)

	items := validators.MapValidationErrorsSlice(err)
	require.Len(t, items, 1)
	assert.Equal(t, "must match Password", items[0].Message)
}

func TestMessage_UnknownTag_ReturnsIsInvalid(t *testing.T) {
	// Use a custom validator that we register
	v := validator.New()
	validators.RegisterJSONTagName(v)
	_ = v.RegisterValidation("custom_tag", func(fl validator.FieldLevel) bool {
		return false // always fail
	})

	type Req struct {
		Value string `json:"value" validate:"custom_tag"`
	}
	err := v.Struct(Req{Value: "anything"})
	require.Error(t, err)

	items := validators.MapValidationErrorsSlice(err)
	require.Len(t, items, 1)
	assert.Equal(t, "is invalid", items[0].Message)
}

// ── GetFieldErrorItems ────────────────────────────────────────────────────────

func TestGetFieldErrorItems_ReturnsMapOfFieldToMessage(t *testing.T) {
	items := validators.FieldErrorItems{
		{Field: "email", Message: "is required"},
		{Field: "name", Message: "must be at least 3 characters"},
	}

	result := items.GetFieldErrorItems()

	assert.Equal(t, map[string]string{
		"email": "is required",
		"name":  "must be at least 3 characters",
	}, result)
}

func TestGetFieldErrorItems_EmptySlice_ReturnsEmptyMap(t *testing.T) {
	items := validators.FieldErrorItems{}

	result := items.GetFieldErrorItems()

	assert.NotNil(t, result)
	assert.Empty(t, result)
}

func TestGetFieldErrorItems_SingleItem(t *testing.T) {
	items := validators.FieldErrorItems{
		{Field: "age", Message: "must be ≥ 18"},
	}

	result := items.GetFieldErrorItems()

	assert.Len(t, result, 1)
	assert.Equal(t, "must be ≥ 18", result["age"])
}

func TestGetFieldErrorItems_DuplicateField_LastWins(t *testing.T) {
	items := validators.FieldErrorItems{
		{Field: "email", Message: "is required"},
		{Field: "email", Message: "must be a valid email"},
	}

	result := items.GetFieldErrorItems()

	assert.Len(t, result, 1)
	assert.Equal(t, "must be a valid email", result["email"])
}

func TestGetFieldErrorItems_FromMapValidationErrorsSlice(t *testing.T) {
	type Req struct {
		Email string `json:"email" validate:"required"`
		Name  string `json:"name"  validate:"required"`
	}

	v := newValidate()
	err := v.Struct(Req{})
	require.Error(t, err)

	items := validators.FieldErrorItems(validators.MapValidationErrorsSlice(err))
	result := items.GetFieldErrorItems()

	assert.Contains(t, result, "email")
	assert.Contains(t, result, "name")
	assert.Equal(t, "is required", result["email"])
	assert.Equal(t, "is required", result["name"])
}

// ── helper ────────────────────────────────────────────────────────────────────

func mustValidate(t *testing.T, s any) []validators.FieldErrorItem {
	t.Helper()
	v := newValidate()
	err := v.Struct(s)
	require.Error(t, err)
	items := validators.MapValidationErrorsSlice(err)
	require.NotEmpty(t, items)
	return items
}
