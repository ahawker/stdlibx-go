package stdlib

import (
	"fmt"
	"reflect"
	"time"
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
		panic(fmt.Sprintf("MustT[%s] received %T value", reflect.TypeFor[T](), v))
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

// MustMapAny returns the map[string]any of the given value and panics if it cannot.
func MustMapAny[T any](value T) map[string]any {
	v, err := ToMapAny[T](value)
	if err != nil {
		panic(err)
	}
	return v
}

// MustMapString returns the map[string]string of the given value and panics if it cannot.
func MustMapString[T any](value T) map[string]string {
	v, err := ToMapString[T](value)
	if err != nil {
		panic(err)
	}
	return v
}

// MustString returns the string representation of the given value and panics if it cannot.
func MustString[T any](value T) string {
	v, err := ToString[T](value)
	if err != nil {
		panic(err)
	}
	return v
}

// MustBool returns the bool representation of the given value and panics if it cannot.
func MustBool[T any](value T) bool {
	v, err := ToBool[T](value)
	if err != nil {
		panic(err)
	}
	return v
}

// MustInt returns the int representation of the given value and panics if it cannot.
func MustInt[T any](value T) int {
	v, err := ToInt[T](value)
	if err != nil {
		panic(err)
	}
	return v
}

// MustDuration returns the time.Duration representation of the given value and panics if it cannot.
func MustDuration[T any](value T) time.Duration {
	v, err := ToDuration[T](value)
	if err != nil {
		panic(err)
	}
	return v
}
