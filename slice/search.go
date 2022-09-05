package slice

import "golang.org/x/exp/constraints"

// Contains returns true if the provided slice contains the target value.
func Contains[T comparable](sl []T, val T) bool {
	for _, v := range sl {
		if v == val {
			return true
		}
	}

	return false
}

// ValueAt returns if exists values[idx] else the fallback value.
func ValueAt[T any](idx int, values []T, fallback T) T {
	if (idx < 0) || (idx >= len(values)) {
		return fallback
	}

	return values[idx]
}

// Min returns the minimum of the provided values.
func Min[T constraints.Ordered](values ...T) T {
	if len(values) == 0 {
		return *new(T)
	}

	res := values[0]
	for _, v := range values[1:] {
		if v < res {
			res = v
		}
	}

	return res
}
