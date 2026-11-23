package pointer_test

import (
	"fmt"

	"github.com/chaihanij/backend-challenge/app/util/pointer"
)

func ExampleToPointer() {
	p := pointer.ToPointer("hello")
	fmt.Println(*p)

	// Zero value (empty string) returns nil.
	n := pointer.ToPointer("")
	fmt.Println(n)
	// Output:
	// hello
	// <nil>
}

func ExampleToValue() {
	p := pointer.ToPointer("hello")
	fmt.Println(pointer.ToValue(p))

	// nil pointer returns the zero value for the type.
	fmt.Println(pointer.ToValue[string](nil))
	// Output:
	// hello
	//
}

func ExampleToPointer_bool() {
	// false is the zero value for bool, so ToPointer returns nil.
	fmt.Println(pointer.ToPointer(false))

	// true is non-zero, so ToPointer returns a pointer.
	p := pointer.ToPointer(true)
	fmt.Println(*p)
	// Output:
	// <nil>
	// true
}
