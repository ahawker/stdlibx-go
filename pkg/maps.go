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
