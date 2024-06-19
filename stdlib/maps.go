package stdlib

// MapFilter will return a new map containing only items
// from the input map that match the predicate function.
func MapFilter[K comparable, V any](input map[K]V, predicate KeyedPredicate[K, V]) map[K]V {
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
func MapFilterRange[K comparable, V any](input KeyedRanger[K, V], predicate KeyedPredicate[K, V]) map[K]V {
	filtered := make(map[K]V)
	input.Range(func(key K, val V) bool {
		if predicate(key, val) {
			filtered[key] = val
		}
		return true
	})
	return filtered
}

// MapKeys returns a slice of all keys for the map.
func MapKeys[K comparable, V any](input map[K]V) []K {
	keys := make([]K, 0, len(input))
	for k := range input {
		keys = append(keys, k)
	}
	return keys
}

// MapValues returns a slice of all values for the map.
func MapValues[K comparable, V any](input map[K]V) []V {
	values := make([]V, 0, len(input))
	for _, v := range input {
		values = append(values, v)
	}
	return values
}
