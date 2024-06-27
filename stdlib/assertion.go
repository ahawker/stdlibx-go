package stdlib

import (
	"reflect"
)

// ErrTypeAssertionFailed is returned when attempting to type assert a value
// that cannot be converted.
var ErrTypeAssertionFailed = Error{
	Code:      "type_assertion_failed",
	Message:   "type assertion failed",
	Namespace: ErrorNamespaceDefault,
}

// As performs a type assertion of 'any' value to type T. If it fails an error is returned.
func As[T any](v any) (T, error) {
	t, ok := v.(T)
	if !ok {
		return *new(T), ErrTypeAssertionFailed.Wrapf("value_type=%T desired_type=%s", v, reflect.TypeFor[T]())
	}
	return t, nil
}
