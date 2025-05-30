package hiter_test

import (
	"cmp"
	"slices"
	"testing"

	"github.com/ngicks/go-iterator-helper/hiter"
	"github.com/ngicks/go-iterator-helper/hiter/mathiter"
	"gotest.tools/v3/assert"
)

func TestMinMax(t *testing.T) {
	type sample struct {
		i int
	}

	rng := hiter.Map(func(i int) int { return i + 5 }, mathiter.Rng(100))

	for i := range 5 {
		nums := slices.Collect(hiter.Limit(2, rng))
		m := hiter.Min(hiter.Concat(hiter.Limit(nums[0], rng), hiter.Once(i), hiter.Limit(nums[1], rng)))
		assert.Equal(t, i, m)
		m = hiter.Max(hiter.Concat(hiter.Limit(nums[0], rng), hiter.Once(i+106), hiter.Limit(nums[1], rng)))
		assert.Equal(t, i+106, m)

		s := hiter.MinFunc(
			func(i, j sample) int {
				return cmp.Compare(i.i, j.i)
			},
			hiter.Map(
				func(i int) sample { return sample{i} },
				hiter.Concat(hiter.Limit(nums[0], rng), hiter.Once(i), hiter.Limit(nums[1], rng)),
			),
		)
		assert.Equal(t, sample{i}, s)
		s = hiter.MaxFunc(
			func(i, j sample) int {
				return cmp.Compare(i.i, j.i)
			},
			hiter.Map(
				func(i int) sample { return sample{i} },
				hiter.Concat(hiter.Limit(nums[0], rng), hiter.Once(i+106), hiter.Limit(nums[1], rng)),
			),
		)
		assert.Equal(t, sample{i + 106}, s)
	}

	assert.Equal(t, 0, hiter.Min(hiter.Empty[int]()))
	assert.Equal(t, 0, hiter.Max(hiter.Empty[int]()))
	assert.Equal(t, sample{0}, hiter.MinFunc(func(i, j sample) int { return cmp.Compare(i.i, j.i) }, hiter.Empty[sample]()))
	assert.Equal(t, sample{0}, hiter.MaxFunc(func(i, j sample) int { return cmp.Compare(i.i, j.i) }, hiter.Empty[sample]()))
}
