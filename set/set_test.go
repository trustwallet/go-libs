package set

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSetFromValues(t *testing.T) {
	s := NewFromValues("foo", "bar")
	assert.Equal(t, 2, s.Size())
	assert.True(t, s.Contains("foo"))
	assert.True(t, s.Contains("bar"))
	assert.False(t, s.Contains("baz"))
}

func TestSet(t *testing.T) {
	s := New[string]()
	assert.Empty(t, s.ToSlice())

	s.Add("foo")
	s.Add("bar")
	assert.Equal(t, 2, s.Size())

	s.Add("bar", "baz")
	assert.True(t, s.Contains("foo"))
	assert.True(t, s.Contains("bar"))
	assert.True(t, s.Contains("baz"))
	assert.Equal(t, 3, s.Size())

	s.Remove("baz")
	assert.Equal(t, 2, s.Size())

	s.Remove("foo")
	assert.Equal(t, 1, s.Size())

	assert.False(t, s.Contains("foo"))
	assert.True(t, s.Contains("bar"))
	assert.Equal(t, []string{"bar"}, s.ToSlice())

	s.ResetFromSlice([]string{"a", "b", "a", "c", "c", "c"})
	assert.Equal(t, 3, s.Size())

	s.Clear()
	assert.False(t, s.Contains("bar"))
	assert.Empty(t, s.ToSlice())
	assert.Equal(t, 0, s.Size())
}
