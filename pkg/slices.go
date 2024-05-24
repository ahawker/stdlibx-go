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
