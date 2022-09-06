package set

import (
	"encoding/json"
)

type Set[T comparable] struct {
	values map[T]struct{}
}

func New[T comparable]() *Set[T] {
	return &Set[T]{
		values: make(map[T]struct{}),
	}
}

func NewFromValues[T comparable](values ...T) *Set[T] {
	s := New[T]()
	s.Add(values...)
	return s
}

func (s *Set[T]) Add(values ...T) {
	for _, val := range values {
		s.values[val] = struct{}{}
	}
}

func (s *Set[T]) Clear() {
	s.values = make(map[T]struct{})
}

func (s *Set[T]) Remove(val T) {
	delete(s.values, val)
}

func (s *Set[T]) Contains(val T) bool {
	_, contains := s.values[val]
	return contains
}

func (s *Set[T]) ContainsAny(values ...T) bool {
	for _, val := range values {
		if s.Contains(val) {
			return true
		}
	}

	return false
}

func (s *Set[T]) ContainsAll(values ...T) bool {
	for _, val := range values {
		if !s.Contains(val) {
			return false
		}
	}

	return true
}

func (s *Set[T]) Extend(s2 *Set[T]) {
	for v := range s2.Values() {
		s.Add(v)
	}
}

func (s *Set[T]) Size() int {
	return len(s.values)
}

func (s *Set[T]) Values() map[T]struct{} {
	return s.values
}

func (s *Set[T]) ToSlice() []T {
	sl := make([]T, 0, len(s.values))

	for v := range s.values {
		sl = append(sl, v)
	}
	return sl
}

func (s *Set[T]) MarshalJSON() ([]byte, error) {
	return json.Marshal(s.ToSlice())
}

func (s *Set[T]) UnmarshalJSON(data []byte) error {
	values := make([]T, 0)

	err := json.Unmarshal(data, &values)
	if err != nil {
		return err
	}

	s.Clear()
	for _, v := range values {
		s.Add(v)
	}

	return nil
}
