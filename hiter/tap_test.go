package hiter_test

import (
	"iter"
	"testing"

	"github.com/ngicks/go-iterator-helper/hiter"
	"gotest.tools/v3/assert"
	"gotest.tools/v3/assert/cmp"
)

func TestTap(t *testing.T) {
	var observed1 []int

	testCase1[int]{
		Seq: func() iter.Seq[int] {
			observed1 = observed1[:0]
			return hiter.Tap(
				func(i int) {
					observed1 = append(observed1, i)
				},
				hiter.Range(0, 5),
			)
		},
		Expected: []int{0, 1, 2, 3, 4},
		BreakAt:  2,
	}.Test(t, func(_, count int) {
		switch count {
		case 0:
			assert.Assert(t, cmp.DeepEqual([]int{0, 1, 2, 3, 4}, observed1))
		case 1:
			assert.Assert(t, cmp.DeepEqual([]int{0, 1, 2}, observed1))
		}
	})

	var observed2 hiter.KeyValues[int, int]

	testCase2[int, int]{
		Seq: func() iter.Seq2[int, int] {
			observed2 = observed2[:0]
			return hiter.Tap2(
				func(i, j int) {
					observed2 = append(observed2, hiter.KeyValue[int, int]{i, j})
				},
				hiter.Pairs(hiter.Range(5, 0), hiter.Range(0, 5)),
			)
		},
		Expected: hiter.KeyValues[int, int]{{5, 0}, {4, 1}, {3, 2}, {2, 3}, {1, 4}},
		BreakAt:  2,
	}.Test(t, func(_, count int) {
		switch count {
		case 0:
			assert.Assert(
				t,
				cmp.DeepEqual(
					hiter.KeyValues[int, int]{{5, 0}, {4, 1}, {3, 2}, {2, 3}, {1, 4}},
					observed2,
				),
			)
		case 1:
			assert.Assert(
				t,
				cmp.DeepEqual(
					hiter.KeyValues[int, int]{{5, 0}, {4, 1}, {3, 2}},
					observed2,
				),
			)
		}
	})

}
