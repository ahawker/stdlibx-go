package stdlibx

// Predicate describes functions which return true/false based on a given input.
type Predicate[T any] func(t T) bool

// KeyedPredicate describes functions which return true/false based on a given
// key/value input.
type KeyedPredicate[K comparable, V any] func(k K, v V) bool
