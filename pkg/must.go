package stdlibx

import "fmt"

// Must panics if given value is equal to the zero value fo the type.
func Must[T any](t T) T {
	if IsZero[T](t) {
		panic(fmt.Sprintf("Must[%T] received zero value", t))
	}
	return t
}
