package stdlib

import "reflect"

// Zeroer describes types that support `IsZero` checks.
type Zeroer interface {
	// IsZero returns true if the value of the instance is equal
	// to the "zero" value of the type.
	IsZero() bool
}

// IsZero returns true if the given value is equal to the
// zero value of the type.
func IsZero[T any](t T) bool {
	switch x := any(t).(type) {
	case Zeroer:
		return x.IsZero()
	default:
		return reflect.DeepEqual(t, *new(T))
	}
}
