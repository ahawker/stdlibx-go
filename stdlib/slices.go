package stdlib

// SliceFlatten will flatten a slice of slices into a
// single slice.
func SliceFlatten[T any](input ...[]T) []T {
	var output []T
	for _, slice := range input {
		output = append(output, slice...)
	}
	return output
}

// SliceSet returns a set from the given slice.
func SliceSet[T comparable](input []T) map[T]struct{} {
	output := make(map[T]struct{}, len(input))
	for _, item := range input {
		output[item] = struct{}{}
	}
	return output
}

// SliceTypeAssert takes a slice of one type and asserts individual
// items to the other.
func SliceTypeAssert[TIn any, TOut any](input []TIn) []TOut {
	output := make([]TOut, 0, len(input))
	for _, item := range input {
		output = append(output, any(item).(TOut))
	}
	return output
}

// SliceMap returns a slice with the results from the given 'map' function.
func SliceMap[TIn any, TOut any](input []TIn, mapper Mapper[TIn, TOut]) []TOut {
	output := make([]TOut, 0, len(input))
	for _, item := range input {
		output = append(output, mapper(item))
	}
	return output
}

// SliceToMap returns a map from the given slice and key function.
func SliceToMap[K comparable, V any](input []V, key func(v V) K) map[K]V {
	output := make(map[K]V, len(input))
	for _, item := range input {
		output[key(item)] = item
	}
	return output
}

// SliceFilter will return a new slice containing only items
// from the given input that match the predicate function.
func SliceFilter[T any](input []T, predicate Predicate[T]) []T {
	var output []T
	for _, item := range input {
		if predicate(item) {
			output = append(output, item)
		}
	}
	return output
}

// SliceFilterRange will return a new slice containing only items
// from the given input ranger that match the predicate function.
func SliceFilterRange[T any](input Ranger[T], predicate Predicate[T]) []T {
	var output []T
	input.Range(func(item T) bool {
		if predicate(item) {
			output = append(output, item)
		}
		return true
	})
	return output
}
