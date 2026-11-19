package pointer_test

import (
	"testing"

	"github.com/chaihanij/backend-challenge/app/util/pointer"

	"github.com/stretchr/testify/assert"
)

func TestToValue_Int_NonNil(t *testing.T) {
	val := 42
	result := pointer.ToValue(&val)

	assert.Equal(t, 42, result)
}

func TestToValue_Int_Nil(t *testing.T) {
	result := pointer.ToValue[int](nil)

	assert.Equal(t, 0, result)
}

func TestToValue_String_NonNil(t *testing.T) {
	val := "hello"
	result := pointer.ToValue(&val)

	assert.Equal(t, "hello", result)
}

func TestToValue_String_Nil(t *testing.T) {
	result := pointer.ToValue[string](nil)

	assert.Equal(t, "", result)
}

func TestToValue_Float_NonNil(t *testing.T) {
	val := 3.14
	result := pointer.ToValue(&val)

	assert.InDelta(t, 3.14, result, 0.0001)
}

func TestToValue_Float_Nil(t *testing.T) {
	result := pointer.ToValue[float64](nil)

	assert.Equal(t, 0.0, result)
}

func TestToValue_Bool_True(t *testing.T) {
	val := true
	result := pointer.ToValue(&val)

	assert.True(t, result)
}

func TestToValue_Bool_Nil(t *testing.T) {
	result := pointer.ToValue[bool](nil)

	assert.False(t, result)
}

func TestToValue_Struct_NonNil(t *testing.T) {
	s := sampleStruct{Name: "Alice", Age: 30}
	result := pointer.ToValue(&s)

	assert.Equal(t, s, result)
}

func TestToValue_Struct_Nil(t *testing.T) {
	result := pointer.ToValue[sampleStruct](nil)

	assert.Equal(t, sampleStruct{}, result)
}

func TestToValue_ReturnsDerefencedValue(t *testing.T) {
	val := 100
	ptr := &val
	result := pointer.ToValue(ptr)

	// แก้ค่าผ่าน pointer ภายหลัง ไม่ควรกระทบ result ที่ได้ไปแล้ว
	*ptr = 999

	assert.Equal(t, 100, result)
}
