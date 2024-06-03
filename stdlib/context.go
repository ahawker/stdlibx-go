package stdlib

import (
	"context"
	"fmt"
	"reflect"
)

// Based on https://github.com/tailscale/tailscale/blob/main/util/ctxkey/key.go

// NewContextKey creates a new context key for a generic type.
func NewContextKey[T any](name string, def T) ContextKey[T] {
	if name == "" {
		name = reflect.TypeFor[T]().String()
	}
	key := ContextKey[T]{name: &stringer[string]{name}}
	if dv := reflect.ValueOf(def); dv.IsValid() && !dv.IsZero() {
		key.def = &def
	}
	return key
}

// ContextKey is a generic key type associated with a specific value type. This
// should be used with non-exported Go types to avoid potential key collisions
// within the context object.
type ContextKey[T any] struct {
	name *stringer[string]
	def  *T
}

// WithValue returns a copy of parent in which the value associated with key is value.
//
// It is a type-safe equivalent of [context.WithValue].
func (k ContextKey[T]) WithValue(parent context.Context, value T) context.Context {
	return context.WithValue(parent, k.name, stringer[T]{value})
}

// ValueOk returns the value in the context associated with this key
// and also reports whether it was present.
// If the value is not present, it returns the default value.
func (k ContextKey[T]) ValueOk(ctx context.Context) (T, bool) {
	cv, ok := ctx.Value(k.name).(stringer[T])
	if !ok && k.def != nil {
		cv.t = *k.def
	}
	return cv.t, ok
}

// Value returns the value in the context associated with this key.
// If the value is not present, it returns the default value.
func (k ContextKey[T]) Value(ctx context.Context) T {
	v, _ := k.ValueOk(ctx)
	return v
}

// String returns the name of the key.
func (k ContextKey[T]) String() string {
	return k.name.String()
}

// stringer supports the 'fmt.Stringer' interface for arbitrary generic types.
type stringer[T any] struct{ t T }

// String returns the name.
func (g stringer[T]) String() string {
	return fmt.Sprint(g.t)
}
