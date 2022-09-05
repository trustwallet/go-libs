package slice

type Unique[T comparable] struct {
	valuesSet map[T]struct{}
	values    []T
}

func NewUnique[T comparable]() *Unique[T] {
	return &Unique[T]{
		valuesSet: make(map[T]struct{}),
		values:    make([]T, 0),
	}
}

func (u *Unique[T]) Add(val T) {
	if _, exists := u.valuesSet[val]; !exists {
		u.valuesSet[val] = struct{}{}
		u.values = append(u.values, val)
	}
}

func (u *Unique[T]) Contains(val T) bool {
	_, contains := u.valuesSet[val]
	return contains
}

func (u *Unique[T]) Values() []T {
	return u.values
}
