package stdlib

// Pointer returns a pointer to the given value.
func Pointer[T any](t T) *T {
	return &t
}

// Dereference dereference the given pointer.
func Dereference[T any](t *T) T {
	return *t
}

// SafeDereference dereference the given pointer and return the zero value if the pointer is nil.
func SafeDereference[T any](t *T) T {
	if t == nil {
		return *new(T)
	}
	return *t
}
