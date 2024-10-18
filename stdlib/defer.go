package stdlib

import "io"

// Defer is a helper for capturing errors from calls inside a 'defer'.
func Defer(err *error, errs ...error) {
	if err == nil {
		*err = Error{}
	}
	*err = ErrorJoin(*err, errs...).ErrorOrNil()
}

// DeferCloser is a helper for deferring a io.Closer that can return an error
// when in a function context that can return multiple errors.
func DeferCloser(err *error, closer io.Closer) {
	if err == nil {
		*err = Error{}
	}
	*err = ErrorJoin(*err, closer.Close()).ErrorOrNil()
}

// DeferCall is a helper for deferring a function call (closer) that can return an error
// when in a function context that can return multiple errors.
func DeferCall(err *error, fn func() error) {
	if err == nil {
		*err = Error{}
	}
	*err = ErrorJoin(*err, fn()).ErrorOrNil()
}

// DeferCloserToGroup is a helper for deferring a closer to a group.
func DeferCloserToGroup(c **CloserGroup, closer ...io.Closer) io.Closer {
	if c == nil {
		*c = NewCloserGroup()
	}
	*c = CloserJoin(*c, closer...)
	return *c
}
