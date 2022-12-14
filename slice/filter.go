package slice

// Filter returns a sub-slice with all the elements that satisfy the fn condition.
func Filter[T any](s []T, fn func(T) bool) []T {
	filtered := make([]T, 0, len(s))

	for _, el := range s {
		if fn(el) {
			filtered = append(filtered, el)
		}
	}

	return filtered
}
