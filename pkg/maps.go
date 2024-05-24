package stdlibx

// MapFilter will return a new map containing only items
// from the input map that match the predicate function.
func MapFilter[K comparable, V any](input map[K]V, predicate func(k K, v V) bool) map[K]V {
	filtered := make(map[K]V)
	for key, val := range input {
		if predicate(key, val) {
			filtered[key] = val
		}
	}
	return filtered
}

// MapFilterRange will return a new map containing only items
// from the input keyed ranger that match the predicate function.
func MapFilterRange[K comparable, V any](input KeyedRanger[K, V], predicate func(k K, v V) bool) map[K]V {
	filtered := make(map[K]V)
	input.Range(func(key K, val V) bool {
		if predicate(key, val) {
			filtered[key] = val
		}
		return true
	})
	return filtered
}
