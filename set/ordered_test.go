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

		assert.Equal(t, 3, u.Size())

		assert.Equal(t, 1, u.ValueAt(0, 0))
		assert.Equal(t, 5, u.ValueAt(1, 0))
		assert.Equal(t, 8, u.ValueAt(2, 0))
		assert.Equal(t, 0, u.ValueAt(256, 0))
		assert.Equal(t, 0, u.ValueAt(-1, 0))
	})

	t.Run("str", func(t *testing.T) {
		u := NewOrderedSet[string]()
		u.Add("foo")
		u.Add("bar")
		u.Add("baz")
		u.Add("foo")
		assert.Equal(t, []string{"foo", "bar", "baz"}, u.Values())

		assert.Equal(t, 3, u.Size())

		assert.Equal(t, "foo", u.ValueAt(0, ""))
		assert.Equal(t, "bar", u.ValueAt(1, ""))
		assert.Equal(t, "baz", u.ValueAt(2, ""))
		assert.Equal(t, "", u.ValueAt(256, ""))
	})
}
