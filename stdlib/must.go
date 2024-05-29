package stdlib

import (
	"errors"
	"fmt"
)

// Must panics if given value is equal to the zero value of the type.
func Must[T any](t T) T {
	if IsZero[T](t) {
		panic(fmt.Sprintf("Must[%T] received zero value", t))
	}
	return t
}

// MustE panics if the given func returns an error for the value returned
// is equal to the zero value of the type.
func MustE[T any](fn func() (T, error)) T {
	t, err := fn()
	if err != nil {
		panic(errors.Join(fmt.Errorf("MustE[%T] func returned error", *new(T)), err))
	}
	return Must[T](t)
}
