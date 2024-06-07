package stdlib

// Valid is a struct that wraps an arbitrary value to indicate that it is
// valid and has passed all checks.
type Valid[T any] struct {
	// Value of type T which has been validated.
	Value T
}

// Validator defines functional validator for type t.
type Validator[T any] func(t T) error

// ValidCheck applies all functional validators to type t
// and returns the error if any fail to apply.
//
// If all return successfully, a Valid[T] is returned.
func ValidCheck[T any](t T, validators ...Validator[T]) (*Valid[T], error) {
	eg := NewErrorGroup()
	for _, v := range validators {
		if err := v(t); err != nil {
			eg.Append(err)
		}
	}
	if err := eg.ErrorOrNil(); err != nil {
		return nil, err
	}
	return &Valid[T]{Value: t}, nil
}
