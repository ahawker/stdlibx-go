package stdlib

// SliceFlatten will flatten a slice of slices into a
// single slice.
func SliceFlatten[T any](slices ...[]T) []T {
	var flattened []T
	for _, slice := range slices {
		flattened = append(flattened, slice...)
	}
	return flattened
}

// SliceToMap returns a map from the given slice and key function.
func SliceToMap[K comparable, V any](input []V, key func(v V) K) map[K]V {
	result := make(map[K]V, len(input))
	for _, item := range input {
		result[key(item)] = item
	}
	return result
}

// SliceFilter will return a new slice containing only items
// from the given input that match the predicate function.
func SliceFilter[T any](input []T, predicate Predicate[T]) []T {
	var filtered []T
	for _, item := range input {
		if predicate(item) {
			filtered = append(filtered, item)
		}
	}
	return filtered
}

// SliceFilterRange will return a new slice containing only items
// from the given input ranger that match the predicate function.
func SliceFilterRange[T any](input Ranger[T], predicate Predicate[T]) []T {
	var filtered []T
	input.Range(func(item T) bool {
		if predicate(item) {
			filtered = append(filtered, item)
		}
		return true
	})
	return filtered
}
