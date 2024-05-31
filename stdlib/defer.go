package stdlib

// Defer is a helper for deferring a function call that can return an error
// when in a function context that can return multiple errors.
func Defer(err *error, fn func() error) {
	if err == nil {
		*err = Error{}
	}
	*err = ErrorJoin(*err, fn()).ErrorOrNil()
}
