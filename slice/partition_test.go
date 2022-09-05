package slice

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPartition(t *testing.T) {
	t.Run("logic", func(t *testing.T) {
		testCasesStr := []struct {
			s     []string
			exp   [][]string
			pSize int
		}{
			{s: []string{"a", "b", "c", "d"}, exp: [][]string{{"a", "b"}, {"c", "d"}}, pSize: 2},
			{s: []string{"a", "b", "c", "d"}, exp: [][]string{{"a", "b", "c", "d"}}, pSize: 8},
			{s: []string{"a", "b", "c", "d"}, exp: [][]string{{"a"}, {"b"}, {"c"}, {"d"}}, pSize: 1},
			{s: []string{"a"}, exp: [][]string{{"a"}}, pSize: 1},
			{s: []string{"a"}, exp: [][]string{{"a"}}, pSize: 100},
			{s: []string{"a"}, exp: [][]string{}, pSize: 0},
			{s: []string{"a", "b", "c", "d", "d", "da"}, exp: [][]string{{"a", "b", "c", "d", "d"}, {"da"}}, pSize: 5},
		}

		for _, tc := range testCasesStr {
			act := Partition(tc.s, tc.pSize)
			assert.Equal(t, tc.exp, act)
		}
	})

	t.Run("generics int", func(t *testing.T) {
		act := Partition([]int{10, 20, 30, 40, 50}, 3)
		assert.Equal(t, [][]int{{10, 20, 30}, {40, 50}}, act)
	})
}
