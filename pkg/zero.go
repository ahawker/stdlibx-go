package stdlibx

import "reflect"

// IsZero returns true if the given value is equal to the
// zero value of the type.
func IsZero[T any](t T) bool {
	return reflect.DeepEqual(t, *new(T))
}
