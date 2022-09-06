package set

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestContains(t *testing.T) {
	t.Run("int", func(t *testing.T) {
		u := NewOrderedSet[int]()
		u.Add(1)
		u.Add(5)
		u.Add(5)
		u.Add(8)
		assert.Equal(t, []int{1, 5, 8}, u.Values())
	})

	t.Run("str", func(t *testing.T) {
		u := NewOrderedSet[string]()
		u.Add("foo")
		u.Add("bar")
		u.Add("baz")
		u.Add("foo")
		assert.Equal(t, []string{"foo", "bar", "baz"}, u.Values())
	})
}
