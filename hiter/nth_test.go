package hiter_test

import (
	"slices"
	"testing"

	"github.com/ngicks/go-iterator-helper/hiter"
	"gotest.tools/v3/assert"
)

func TestNth(t *testing.T) {

	t.Run("1", func(t *testing.T) {
		src := []int{1, 2, 3, 4, 5}
		var (
			v  int
			ok bool
		)
		v, ok = hiter.Nth(-1, slices.Values(src))
		assert.Assert(t, !ok)
		assert.Equal(t, 0, v)

		v, ok = hiter.Nth(0, slices.Values(src))
		assert.Assert(t, ok)
		assert.Equal(t, 1, v)

		v, ok = hiter.Nth(2, slices.Values(src))
		assert.Assert(t, ok)
		assert.Equal(t, 3, v)

		v, ok = hiter.Nth(len(src), slices.Values(src))
		assert.Assert(t, !ok)
		assert.Equal(t, 0, v)
	})
	t.Run("2", func(t *testing.T) {
		src := []int{1, 2, 3, 4, 5}
		var (
			k, v int
			ok   bool
		)
		k, v, ok = hiter.Nth2(-1, slices.All(src))
		assert.Assert(t, !ok)
		assert.Equal(t, 0, k)
		assert.Equal(t, 0, v)

		k, v, ok = hiter.Nth2(0, slices.All(src))
		assert.Assert(t, ok)
		assert.Equal(t, 0, k)
		assert.Equal(t, 1, v)

		k, v, ok = hiter.Nth2(2, slices.All(src))
		assert.Assert(t, ok)
		assert.Equal(t, 2, k)
		assert.Equal(t, 3, v)

		k, v, ok = hiter.Nth2(len(src), slices.All(src))
		assert.Assert(t, !ok)
		assert.Equal(t, 0, k)
		assert.Equal(t, 0, v)
	})
}
