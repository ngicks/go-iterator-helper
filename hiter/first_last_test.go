package hiter_test

import (
	"testing"

	"github.com/ngicks/go-iterator-helper/hiter"
	"gotest.tools/v3/assert"
)

func TestFirstLast(t *testing.T) {
	t.Run("1", func(t *testing.T) {
		src := hiter.Range(0, 5)
		empty := hiter.Range(0, 0)

		var (
			v  int
			ok bool
		)

		v, ok = hiter.First(src)
		assert.Assert(t, ok)
		assert.Equal(t, 0, v)

		v, ok = hiter.First(empty)
		assert.Assert(t, !ok)
		assert.Equal(t, 0, v)

		v, ok = hiter.Last(src)
		assert.Assert(t, ok)
		assert.Equal(t, 4, v)

		v, ok = hiter.Last(empty)
		assert.Assert(t, !ok)
		assert.Equal(t, 0, v)
	})

	t.Run("2", func(t *testing.T) {
		src := hiter.Pairs(hiter.Range(4, -1), hiter.Range(0, 5))
		empty := hiter.Enumerate(hiter.Range(0, 0))

		var (
			k  int
			v  int
			ok bool
		)

		k, v, ok = hiter.First2(src)
		assert.Assert(t, ok)
		assert.Equal(t, 4, k)
		assert.Equal(t, 0, v)

		k, v, ok = hiter.First2(empty)
		assert.Assert(t, !ok)
		assert.Equal(t, 0, k)
		assert.Equal(t, 0, v)

		k, v, ok = hiter.Last2(src)
		assert.Assert(t, ok)
		assert.Equal(t, 0, k)
		assert.Equal(t, 4, v)

		k, v, ok = hiter.Last2(empty)
		assert.Assert(t, !ok)
		assert.Equal(t, 0, k)
		assert.Equal(t, 0, v)
	})
}
