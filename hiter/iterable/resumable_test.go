package iterable_test

import (
	"slices"
	"testing"

	"github.com/ngicks/go-iterator-helper/hiter"
	"github.com/ngicks/go-iterator-helper/hiter/iterable"
	"github.com/ngicks/go-iterator-helper/x/exp/xiter"
	"gotest.tools/v3/assert"
	"gotest.tools/v3/assert/cmp"
)

func TestResumable(t *testing.T) {
	src := hiter.Range(0, 5)

	t.Run("one-run", func(t *testing.T) {
		res := iterable.NewResumable(src)
		assert.Assert(t, cmp.DeepEqual(slices.Collect(src), slices.Collect(res.IntoIter())))
		assert.Assert(t, slices.Collect(res.IntoIter()) == nil)
	})
	t.Run("break-and-resume", func(t *testing.T) {
		res := iterable.NewResumable(src)
		assert.Assert(t, cmp.DeepEqual(slices.Collect(xiter.Limit(src, 3)), slices.Collect(xiter.Limit(res.IntoIter(), 3))))
		assert.Assert(t, cmp.DeepEqual(slices.Collect(hiter.Skip(3, src)), slices.Collect(res.IntoIter())))
		assert.Assert(t, slices.Collect(res.IntoIter()) == nil)
	})
	t.Run("stop", func(t *testing.T) {
		res := iterable.NewResumable(src)
		var values []int
		for i, v := range hiter.Enumerate(res.IntoIter()) {
			values = append(values, v)
			if i == 2 {
				res.Stop()
			}
		}
		assert.Assert(t, cmp.DeepEqual(slices.Collect(xiter.Limit(src, 3)), values))
		assert.Assert(t, slices.Collect(res.IntoIter()) == nil)
	})
}

func TestResumable2(t *testing.T) {
	src := hiter.Enumerate(hiter.Range(0, 5))

	t.Run("one-run", func(t *testing.T) {
		res := iterable.NewResumable2(src)
		defer res.Stop()
		assert.Assert(t, cmp.DeepEqual(hiter.Collect2(src), hiter.Collect2(res.IntoIter2())))
		assert.Assert(t, hiter.Collect2(res.IntoIter2()) == nil)
	})
	t.Run("break-and-resume", func(t *testing.T) {
		res := iterable.NewResumable2(src)
		defer res.Stop()
		assert.Assert(t, cmp.DeepEqual(hiter.Collect2(xiter.Limit2(src, 3)), hiter.Collect2(xiter.Limit2(res.IntoIter2(), 3))))
		assert.Assert(t, cmp.DeepEqual(hiter.Collect2(hiter.Skip2(3, src)), hiter.Collect2(res.IntoIter2())))
		assert.Assert(t, hiter.Collect2(res.IntoIter2()) == nil)
	})
	t.Run("stop", func(t *testing.T) {
		res := iterable.NewResumable2(src)
		defer res.Stop()
		var pairs []hiter.KeyValue[int, int]
		for k, v := range res.IntoIter2() {
			pairs = append(pairs, hiter.KeyValue[int, int]{K: k, V: v})
			if k == 2 {
				res.Stop()
			}
		}
		assert.Assert(t, cmp.DeepEqual(hiter.Collect2(xiter.Limit2(src, 3)), pairs))
		assert.Assert(t, hiter.Collect2(res.IntoIter2()) == nil)
	})
}
