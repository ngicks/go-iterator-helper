package collection

import (
	"iter"
	"maps"
	"slices"
	"testing"

	"github.com/ngicks/go-iterator-helper/hiter"
	"gotest.tools/v3/assert"
	"gotest.tools/v3/assert/cmp"
)

// avoiding xiter dependency.
func mapIter[V1, V2 any](seq iter.Seq[V1], mapper func(V1) V2) iter.Seq[V2] {
	return func(yield func(V2) bool) {
		for v := range seq {
			if !yield(mapper(v)) {
				return
			}
		}
	}
}

func TestPermutations(t *testing.T) {
	m := maps.Collect(
		hiter.Pairs(
			mapIter(
				Permutations([]int{1, 2, 3, 4, 5}),
				func(v []int) [4]int { return [4]int(slices.Clone(v)) },
			),
			hiter.Repeat(struct{}{}, -1),
		),
	)
	assert.Assert(t, cmp.Len(m, 5*4*3*2*1))
}
