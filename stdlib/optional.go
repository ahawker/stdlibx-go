package stdlib

// Some creates an Optional with a value.
func Some[T any](v T) Optional[T] {
	return Optional[T]{value: v, def: *new(T), changed: true}
}

// Default creates an Optional with an empty value and distinct default value.
func Default[T any](v T) Optional[T] {
	return Optional[T]{value: *new(T), def: v, changed: false}
}

// Optional wraps a value of type `T` tracks a default value
// and whether it changed.
type Optional[T any] struct {
	// value that has been set from a call to 'Set'.
	value T
	// def is the default value if not set.
	def T
	// changed is true if the value has been set.
	changed bool
}

// Changed returns true if the value has been set.
func (o *Optional[T]) Changed() bool { return o.changed }

// Get returns the value if it has been set, otherwise the default value.
func (o *Optional[T]) Get() T {
	if o.changed {
		return o.value
	} else {
		return o.def
	}
}

// Set sets the value and marks it as changed.
func (o *Optional[T]) Set(v T) {
	o.value = v
	o.changed = true
}
