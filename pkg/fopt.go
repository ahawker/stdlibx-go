package stdlibx

// FOpt defines functional options for type t.
type FOpt[T any] func(t T) error

// FOptApply applies all functional options to type t
// and returns the error if any fail to apply.
func FOptApply[T any](t T, options ...FOpt[T]) (T, error) {
	for _, o := range options {
		if err := o(t); err != nil {
			return t, err
		}
	}
	return t, nil
}
