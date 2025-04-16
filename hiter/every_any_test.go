package hiter_test

import (
	"testing"

	"github.com/ngicks/go-iterator-helper/hiter"
	"github.com/ngicks/go-iterator-helper/hiter/internal/adapter"
	"gotest.tools/v3/assert"
)

func TestEveryAny(t *testing.T) {
	even := adapter.Map(func(i int) int { return i * 2 }, hiter.Range(0, 10))
	even2 := hiter.Divide(func(i int) (int, int) { return i, i }, even)

	isEven := func(i int) bool { return i%2 == 0 }
	isEven2 := func(i, j int) bool { return i%2 == 0 && j%2 == 0 }
	isOdd := func(i int) bool { return i%2 != 0 }
	isOdd2 := func(i, j int) bool { return i%2 != 0 && j%2 != 0 }

	t.Run("Every", func(t *testing.T) {
		assert.Assert(t, hiter.Every(isEven, hiter.Empty[int]()))
		assert.Assert(t, hiter.Every2(isEven2, hiter.Empty2[int, int]()))

		assert.Assert(t, hiter.Every(isEven, even))
		assert.Assert(t, hiter.Every2(isEven2, even2))

		assert.Assert(t, !hiter.Every(isEven, adapter.Concat(even, hiter.Once(5))))
		assert.Assert(t, !hiter.Every2(isEven2, adapter.Concat2(even2, hiter.Once2(5, 2))))
	})

	t.Run("Any", func(t *testing.T) {
		assert.Assert(t, !hiter.Any(isEven, hiter.Empty[int]()))
		assert.Assert(t, !hiter.Any2(isEven2, hiter.Empty2[int, int]()))

		assert.Assert(t, hiter.Any(isEven, even))
		assert.Assert(t, hiter.Any2(isEven2, even2))

		assert.Assert(t, !hiter.Any(isOdd, even))
		assert.Assert(t, !hiter.Any2(isOdd2, even2))

		assert.Assert(t, hiter.Any(isEven, adapter.Concat(hiter.Once(7), even, hiter.Once(5))))
		assert.Assert(t, hiter.Any2(isEven2, adapter.Concat2(hiter.Once2(7, 8), even2, hiter.Once2(5, 2))))

		assert.Assert(t, hiter.Any(isOdd, adapter.Concat(even, hiter.Once(5))))
		assert.Assert(t, hiter.Any2(isOdd2, adapter.Concat2(even2, hiter.Once2(5, 5))))
	})
}
