package stdlibx

import "fmt"

// Must panics if given value is nil.
func Must[T any](t T) T {
	if t == nil {
		panic(fmt.Sprintf("Must[%T] received nil value", t))
	}
	return t
}

// MustNonZero panics if given value is "zero" value for type t.
func MustNonZero[T any](t T) T {
	if IsZero[T](t) {
		panic(fmt.Sprintf("Must[%T] received zero value", t))
	}
	return t
}
