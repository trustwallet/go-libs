package slice

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSliceContains(t *testing.T) {
	assert.True(t, Contains([]int{1, 5, 100, 1000}, 100))
	assert.True(t, Contains([]string{"abc", "z", "fe", "ll"}, "fe"))
	assert.False(t, Contains([]string{"abc", "z", "fe", "ll"}, ""))
	assert.False(t, Contains([]int{1, 5, 100, 1000}, -1))
	assert.False(t, Contains([]bool{false, false, false}, true))
}

func TestValueAt(t *testing.T) {
	t.Run("logic", func(t *testing.T) {
		type args struct {
			idx      int
			values   []string
			fallback string
		}

		tests := []struct {
			name string
			args args
			want string
		}{
			{
				name: "empty slice negative index",
				args: args{idx: -1, values: nil, fallback: ""},
				want: "",
			},
			{
				name: "empty slice zero index",
				args: args{idx: 0, values: nil, fallback: ""},
				want: "",
			},
			{
				name: "index above bounds",
				args: args{idx: 4, values: []string{"a", "b", "c", "d"}, fallback: "not found"},
				want: "not found",
			},
			{
				name: "negative index",
				args: args{idx: -1, values: []string{"a", "b", "c", "d"}, fallback: "not found"},
				want: "not found",
			},
			{
				name: "zero index",
				args: args{idx: 0, values: []string{"a", "b", "c", "d"}, fallback: "not found"},
				want: "a",
			},
			{
				name: "fourth element",
				args: args{idx: 3, values: []string{"a", "b", "c", "d"}, fallback: "not found"},
				want: "d",
			},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				if got := ValueAt(tt.args.idx, tt.args.values, tt.args.fallback); got != tt.want {
					t.Errorf("StrSliceValueAt() = %v, want %v", got, tt.want)
				}
			})
		}
	})

	t.Run("generics", func(t *testing.T) {
		assert.Equal(t, "foo", ValueAt(5, []string{"a", "B"}, "foo"))
		assert.Equal(t, "B", ValueAt(1, []string{"a", "B", "c"}, "foo"))
		assert.Equal(t, -1, ValueAt(1111, []int{10, 20, 30}, -1))
		assert.Equal(t, 10, ValueAt(0, []int{10, 20, 30}, -1))
	})
}

func TestMin(t *testing.T) {
	assert.Equal(t, 5, Min(10, 30, 5, 123, 99))
	assert.Equal(t, "b", Min("z", "b", "d"))
}
