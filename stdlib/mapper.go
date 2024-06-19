package stdlib

// Mapper describes a 'map' function applied to another type.
type Mapper[TIn any, TOut any] func(t TIn) TOut

// KeyedMapper describes a 'map' function for a given
// key/value input.
type KeyedMapper[K comparable, VIn any, VOut any] func(k K, v VIn) VOut
