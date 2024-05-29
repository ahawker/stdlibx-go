package stdlib

// Ranger describes types that export a `Range` method for iteration
// over single item collections.
type Ranger[T any] interface {
	// Range calls the given function for all items available for iteration.
	//
	// If Range returns `false`, iteration will stop.
	Range(predicate Predicate[T])
}

// KeyedRanger describes types that export a `Range` method for iteration
// over key/value collections.
type KeyedRanger[K comparable, V any] interface {
	// Range calls the given function for all items available for iteration.
	//
	// If Range returns `false`, iteration will stop.
	Range(predicate KeyedPredicate[K, V])
}
