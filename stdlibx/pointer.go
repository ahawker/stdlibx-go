package stdlibx

// Pointer returns a pointer to the given value.
func Pointer[T any](t T) *T {
	return &t
}
