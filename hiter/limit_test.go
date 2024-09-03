package hiter_test

import (
	"iter"
	"slices"
	"testing"

	"github.com/ngicks/go-iterator-helper/hiter"
)

func TestLimitUntil(t *testing.T) {
	src := []int{1, 1, 4, 2, 1, 1, 4, 2, 1, 1, 6, 2, 1, 1, 9, 2}
	t.Run("LimitUntil", func(t *testing.T) {
		testCase1[int]{
			Seq: func() iter.Seq[int] {
				return hiter.LimitUntil(slices.Values(src), func(e int) bool { return e < 6 })
			},
			Expected: []int{1, 1, 4, 2, 1, 1, 4, 2, 1, 1},
			BreakAt:  3,
		}.Test(t)
	})

	t.Run("LimitUntil2", func(t *testing.T) {
		testCase2[int, int]{
			Seq: func() iter.Seq2[int, int] {
				return hiter.LimitUntil2(slices.All(src), func(i, e int) bool { return e < 6 })
			},
			Expected: []hiter.KeyValue[int, int]{{0, 1}, {1, 1}, {2, 4}, {3, 2}, {4, 1}, {5, 1}, {6, 4}, {7, 2}, {8, 1}, {9, 1}},
			BreakAt:  3,
		}.Test(t)
	})

	t.Run("LimitUntil2", func(t *testing.T) {
		testCase2[int, int]{
			Seq: func() iter.Seq2[int, int] {
				return hiter.LimitUntil2(slices.All(src), func(i, e int) bool { return i < 3 })
			},
			Expected: []hiter.KeyValue[int, int]{{0, 1}, {1, 1}, {2, 4}},
			BreakAt:  2,
		}.Test(t)
	})
}
