package stdlib

import (
	"fmt"
	"reflect"
)

// Must panics if given value is equal to the zero value of the type.
func Must[T any](t T) T {
	if IsZero[T](t) {
		panic(fmt.Sprintf("Must[%T] received zero value", t))
	}
	return t
}

// MustT panics if given 'any' value cannot be aliased to the type t.
func MustT[T any](v any) T {
	t, ok := v.(T)
	if !ok {
		panic(fmt.Sprintf("MustT[%T] received %T value", reflect.TypeFor[T](), v))
	}
	return t
}

// MustE panics if the given func returns an error for the value returned
// is equal to the zero value of the type.
func MustE[T any](fn func() (T, error)) T {
	t, err := fn()
	if err != nil {
		panic(ErrorJoin(fmt.Errorf("MustE[%s] func returned error", reflect.TypeFor[T]()), err))
	}
	return Must[T](t)
}
