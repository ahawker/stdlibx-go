package stdlib

import (
	"fmt"
	"reflect"
	"strconv"
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

// MustString returns the string representation of the given value and panics if it cannot.
func MustString[T any](value T) string {
	switch v := any(value).(type) {
	case string:
		return v
	default:
		return fmt.Sprintf("%v", v)
	}
}

// MustBool returns the bool representation of the given value and panics if it cannot.
func MustBool[T any](value T) bool {
	switch v := any(value).(type) {
	case bool:
		return v
	default:
		return !IsZero[T](value)
	}
}

// MustInt returns the int representation of the given value and panics if it cannot.
func MustInt[T any](value T) int {
	switch v := any(value).(type) {
	case int:
		return v
	case float32:
		return int(v)
	case float64:
		return int(v)
	case int32:
		return int(v)
	case int64:
		return int(v)
	case int8:
		return int(v)
	case int16:
		return int(v)
	case uint:
		return int(v)
	case uint8:
		return int(v)
	case uint16:
		return int(v)
	case string:
		i, err := strconv.ParseInt(v, 10, 64)
		if err != nil {
			panic(err)
		}
		return int(i)
	default:
		panic(fmt.Sprintf("cannot convert %T to int", value))
	}
}

// MustDuration returns the time.Duration representation of the given value and panics if it cannot.
func MustDuration[T any](value T) time.Duration {
	switch v := any(value).(type) {
	case time.Duration:
		return v
	case string:
		d, err := time.ParseDuration(v)
		if err != nil {
			panic(err)
		}
		return d
	default:
		panic(fmt.Sprintf("cannot convert %T to time.Duration", value))
	}
}
