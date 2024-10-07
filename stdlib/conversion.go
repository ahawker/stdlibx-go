package stdlib

import (
	"fmt"
	"reflect"
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

// ToMapAny returns the map[string]any representation of the given value and errors if it cannot.
func ToMapAny[T any](value T) (map[string]any, error) {
	switch v := any(value).(type) {
	case map[string]any:
		return v, nil
	default:
		rv := reflect.ValueOf(value)
		switch rv.Kind() {
		case reflect.Map:
			result := make(map[string]any, rv.Len())
			for _, mk := range rv.MapKeys() {
				mapKey, err := ToString(mk.Interface())
				if err != nil {
					return nil, ErrTypeConversionFailed.Wrapf("value_type=%T desired_type=string", mk.Interface())
				}
				result[mapKey] = rv.MapIndex(mk).Interface()
			}
			return result, nil
		default:
			return nil, ErrTypeConversionFailed.Wrapf("value_type=%T desired_type=map[string]any", value)
		}
	}
}

// ToMapString returns the map[string]any representation of the given value and errors if it cannot.
func ToMapString[T any](value T) (map[string]string, error) {
	switch v := any(value).(type) {
	case map[string]string:
		return v, nil
	default:
		rv := reflect.ValueOf(value)
		switch rv.Kind() {
		case reflect.Map:
			result := make(map[string]string, rv.Len())
			for _, mk := range rv.MapKeys() {
				mapKey, err := ToString(mk.Interface())
				if err != nil {
					return nil, ErrTypeConversionFailed.Wrapf("value_type=%T desired_type=string", mk.Interface())
				}
				mapVal, err := ToString(rv.MapIndex(mk).Interface())
				if err != nil {
					return nil, ErrTypeConversionFailed.Wrapf("value_type=%T desired_type=string", mk.Interface())
				}
				result[mapKey] = mapVal
			}
			return result, nil
		default:
			return nil, ErrTypeConversionFailed.Wrapf("value_type=%T desired_type=map[string]any", value)
		}
	}
}

// ToString returns the string representation of the given value and errors if it cannot.
func ToString[T any](value T) (string, error) {
	switch v := any(value).(type) {
	case string:
		return v, nil
	case fmt.Stringer:
		return v.String(), nil
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
