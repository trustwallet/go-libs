package set

type OrderedSet[T comparable] struct {
	valuesSet map[T]struct{}
	values    []T
}

func NewOrderedSet[T comparable]() *OrderedSet[T] {
	return &OrderedSet[T]{
		valuesSet: make(map[T]struct{}),
		values:    make([]T, 0),
	}
}

func (u *OrderedSet[T]) Add(val T) {
	if _, exists := u.valuesSet[val]; !exists {
		u.valuesSet[val] = struct{}{}
		u.values = append(u.values, val)
	}
}

func (u *OrderedSet[T]) Contains(val T) bool {
	_, contains := u.valuesSet[val]
	return contains
}

func (u *OrderedSet[T]) Values() []T {
	return u.values
}
