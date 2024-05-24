package stdlibx

// SliceFlatten will flatten a slice of slices into a
// single slice.
func SliceFlatten[T any](slices ...[]T) ([]T, error) {
	var flattened []T
	for _, slice := range slices {
		flattened = append(flattened, slice...)
	}
	return flattened, nil
}

// SliceFilter will return a new slice containing only items
// from the given input that match the predicate function.
func SliceFilter[T any](input []T, predicate func(t T) bool) []T {
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
func SliceFilterRange[T any](input Ranger[T], predicate func(t T) bool) []T {
	var filtered []T
	input.Range(func(item T) bool {
		if predicate(item) {
			filtered = append(filtered, item)
		}
		return true
	})
	return filtered
}
