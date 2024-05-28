package stdlibx

// Option defines functional options for type t.
type Option[T any] func(t T) error

// OptionApply applies all functional options to type t
// and returns the error if any fail to apply.
func OptionApply[T any](t T, options ...Option[T]) (T, error) {
	for _, o := range options {
		if err := o(t); err != nil {
			return t, err
		}
	}
	return t, nil
}
