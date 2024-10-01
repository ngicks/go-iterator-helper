package hiter_test

import (
	"cmp"
	"slices"
	"testing"

	"github.com/ngicks/go-iterator-helper/hiter"
	"github.com/ngicks/go-iterator-helper/hiter/sh"
	"github.com/ngicks/go-iterator-helper/x/exp/xiter"
	"gotest.tools/v3/assert"
)

func TestMinMax(t *testing.T) {
	type sample struct {
		i int
	}

	rng := xiter.Map(func(i int) int { return i + 5 }, sh.Rng(100))

	for i := range 5 {
		nums := slices.Collect(xiter.Limit(rng, 2))
		m := hiter.Min(xiter.Concat(xiter.Limit(rng, nums[0]), hiter.Once(i), xiter.Limit(rng, nums[1])))
		assert.Equal(t, i, m)
		m = hiter.Max(xiter.Concat(xiter.Limit(rng, nums[0]), hiter.Once(i+106), xiter.Limit(rng, nums[1])))
		assert.Equal(t, i+106, m)

		s := hiter.MinFunc(
			func(i, j sample) int {
				return cmp.Compare(i.i, j.i)
			},
			xiter.Map(
				func(i int) sample { return sample{i} },
				xiter.Concat(xiter.Limit(rng, nums[0]), hiter.Once(i), xiter.Limit(rng, nums[1])),
			),
		)
		assert.Equal(t, sample{i}, s)
		s = hiter.MaxFunc(
			func(i, j sample) int {
				return cmp.Compare(i.i, j.i)
			},
			xiter.Map(
				func(i int) sample { return sample{i} },
				xiter.Concat(xiter.Limit(rng, nums[0]), hiter.Once(i+106), xiter.Limit(rng, nums[1])),
			),
		)
		assert.Equal(t, sample{i + 106}, s)
	}

	assert.Equal(t, 0, hiter.Min(hiter.Empty[int]()))
	assert.Equal(t, 0, hiter.Max(hiter.Empty[int]()))
	assert.Equal(t, sample{0}, hiter.MinFunc(func(i, j sample) int { return cmp.Compare(i.i, j.i) }, hiter.Empty[sample]()))
	assert.Equal(t, sample{0}, hiter.MaxFunc(func(i, j sample) int { return cmp.Compare(i.i, j.i) }, hiter.Empty[sample]()))
}
