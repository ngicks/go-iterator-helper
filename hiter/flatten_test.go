package hiter_test

import (
	"iter"
	"slices"
	"testing"

	"github.com/ngicks/go-iterator-helper/hiter"
)

func TestFlatten(t *testing.T) {
	t.Run("Flatten", func(t *testing.T) {
		testCase1[int]{
			Seq: func() iter.Seq[int] {
				return hiter.Flatten(slices.Values([][]int{{0, 1, 2}, {3, 4}, {5}}))
			},
			Expected: []int{0, 1, 2, 3, 4, 5},
			BreakAt:  4,
		}.Test(t)
	})
	t.Run("FlattenF", func(t *testing.T) {
		testCase2[int, int]{
			Seq: func() iter.Seq2[int, int] {
				return hiter.FlattenF(hiter.Combine(
					slices.Values([][]int{{0, 1, 2}, {3, 4}, {5}}),
					slices.Values([]int{10, 7, 8, 1}),
				))
			},
			Expected: []hiter.KeyValue[int, int]{{0, 10}, {1, 10}, {2, 10}, {3, 7}, {4, 7}, {5, 8}},
			BreakAt:  4,
		}.Test(t)
	})
	t.Run("FlattenL", func(t *testing.T) {
		testCase2[int, int]{
			Seq: func() iter.Seq2[int, int] {
				return hiter.FlattenL(hiter.Combine(
					slices.Values([]int{0, 1, 2, 3, 4, 5}),
					slices.Values([][]int{{10}, {7, 8}, {1}}),
				))
			},
			Expected: []hiter.KeyValue[int, int]{{0, 10}, {1, 7}, {1, 8}, {2, 1}},
			BreakAt:  2,
		}.Test(t)
	})
}
