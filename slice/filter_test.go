package slice

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFilter(t *testing.T) {
	t.Run("logic", func(t *testing.T) {
		assert.Equal(t,
			[]string{"a", "c"},
			Filter(
				[]string{"a", "b", "c"},
				func(s string) bool { return s != "b" },
			))
		assert.Equal(t,
			[]int{1, 1},
			Filter(
				[]int{1, 10, 100, 1000, 100, 10, 1},
				func(i int) bool { return i < 10 },
			))

		type item struct {
			price      int
			condition  string
			onDiscount bool
		}

		assert.Len(t,
			Filter(
				[]item{
					{
						price:      115,
						condition:  "new",
						onDiscount: true,
					},
					{
						price:      225,
						condition:  "used",
						onDiscount: false,
					},
					{
						price:      335,
						condition:  "mint",
						onDiscount: true,
					},
				},
				func(i item) bool { return i.onDiscount && i.condition != "used" && i.price < 300 },
			), 1,
		)
	})
}
