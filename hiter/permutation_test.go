package hiter

import (
	"maps"
	"slices"
	"testing"

	"github.com/ngicks/go-iterator-helper/hiter/internal/adapter"
	"gotest.tools/v3/assert"
	"gotest.tools/v3/assert/cmp"
)

func TestPermutations(t *testing.T) {
	m := maps.Collect(
		Pairs(
			adapter.MapIter(
				func(v []int) [5]int { return [5]int(slices.Clone(v)) },
				Permutations([]int{1, 2, 3, 4, 5}),
			),
			Repeat(struct{}{}, -1),
		),
	)
	assert.Assert(t, cmp.Len(m, 5*4*3*2*1))
}
