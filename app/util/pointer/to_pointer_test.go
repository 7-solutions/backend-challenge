package pointer_test

import (
	"testing"

	"github.com/chaihanij/backend-challenge/app/util/pointer"

	"github.com/stretchr/testify/assert"
)

func TestToPointer_Int_NonZero(t *testing.T) {
	result := pointer.ToPointer(42)

	assert.NotNil(t, result)
	assert.Equal(t, 42, *result)
}

func TestToPointer_Int_Zero(t *testing.T) {
	result := pointer.ToPointer(0)

	assert.Nil(t, result)
}

func TestToPointer_String_NonEmpty(t *testing.T) {
	result := pointer.ToPointer("hello")

	assert.NotNil(t, result)
	assert.Equal(t, "hello", *result)
}

func TestToPointer_String_Empty(t *testing.T) {
	result := pointer.ToPointer("")

	assert.Nil(t, result)
}

func TestToPointer_Float_NonZero(t *testing.T) {
	result := pointer.ToPointer(3.14)

	assert.NotNil(t, result)
	assert.InDelta(t, 3.14, *result, 0.0001)
}

func TestToPointer_Float_Zero(t *testing.T) {
	result := pointer.ToPointer(0.0)

	assert.Nil(t, result)
}

func TestToPointer_Bool_True(t *testing.T) {
	result := pointer.ToPointer(true)

	assert.NotNil(t, result)
	assert.True(t, *result)
}

func TestToPointer_Bool_False(t *testing.T) {
	// false คือ zero value ของ bool
	result := pointer.ToPointer(false)

	assert.Nil(t, result)
}

type sampleStruct struct {
	Name string
	Age  int
}

func TestToPointer_Struct_NonZero(t *testing.T) {
	s := sampleStruct{Name: "Alice", Age: 30}
	result := pointer.ToPointer(s)

	assert.NotNil(t, result)
	assert.Equal(t, s, *result)
}

func TestToPointer_Struct_Zero(t *testing.T) {
	result := pointer.ToPointer(sampleStruct{})

	assert.Nil(t, result)
}

func TestToPointer_ReturnsIndependentPointer(t *testing.T) {
	original := 100
	result := pointer.ToPointer(original)

	// แก้ค่าต้นทาง ไม่ควรกระทบ pointer
	original = 999

	assert.Equal(t, 100, *result)
}
