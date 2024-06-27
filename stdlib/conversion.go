package stdlib

import (
	"fmt"
	"strconv"
	"time"
)

// ErrTypeConversionFailed is returned when attempting a type conversion that
// cannot be performed.
var ErrTypeConversionFailed = Error{
	Code:      "type_conversion_failed",
	Message:   "type conversion failed",
	Namespace: ErrorNamespaceDefault,
}

// ToString returns the string representation of the given value and errors if it cannot.
func ToString[T any](value T) (string, error) {
	switch v := any(value).(type) {
	case string:
		return v, nil
	default:
		return fmt.Sprintf("%v", v), nil
	}
}

// ToBool returns the bool representation of the given value and errors if it cannot.
func ToBool[T any](value T) (bool, error) {
	switch v := any(value).(type) {
	case bool:
		return v, nil
	default:
		return !IsZero[T](value), nil
	}
}

// ToInt returns the int representation of the given value and errors if it cannot.
func ToInt[T any](value T) (int, error) {
	switch v := any(value).(type) {
	case int:
		return v, nil
	case float32:
		return int(v), nil
	case float64:
		return int(v), nil
	case int32:
		return int(v), nil
	case int64:
		return int(v), nil
	case int8:
		return int(v), nil
	case int16:
		return int(v), nil
	case uint:
		return int(v), nil
	case uint8:
		return int(v), nil
	case uint16:
		return int(v), nil
	case string:
		i, err := strconv.ParseInt(v, 10, 64)
		if err != nil {
			return 0, ErrTypeConversionFailed.Wrapf("value_type=%T desired_type=int", value).Wrap(err)
		}
		return int(i), nil
	default:
		return 0, ErrTypeConversionFailed.Wrapf("value_type=%T desired_type=int", value)
	}
}

// ToDuration returns the time.Duration representation of the given value and errors if it cannot.
func ToDuration[T any](value T) (time.Duration, error) {
	switch v := any(value).(type) {
	case time.Duration:
		return v, nil
	case string:
		d, err := time.ParseDuration(v)
		if err != nil {
			return 0, ErrTypeConversionFailed.Wrapf("value_type=%T desired_type=time.Duration", value).Wrap(err)
		}
		return d, nil
	default:
		return 0, ErrTypeConversionFailed.Wrapf("value_type=%T desired_type=time.Duration", value)
	}
}
