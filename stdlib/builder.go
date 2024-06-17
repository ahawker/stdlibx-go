package stdlib

// Builder defines functional builder for type t.
type Builder[T any] func(t T) (T, error)

// Build applies all functional builders to type t
// and returns the type built.
func Build[T any](t T, builders ...Builder[T]) (T, error) {
	var err error
	for _, b := range builders {
		t, err = b(t)
		if err != nil {
			return t, err
		}
	}
	return t, nil
}
