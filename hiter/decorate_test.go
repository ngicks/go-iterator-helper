package hiter_test

import (
	"iter"
	"slices"
	"testing"

	"github.com/ngicks/go-iterator-helper/hiter"
	"github.com/ngicks/go-iterator-helper/hiter/iterable"
)

func TestDecorate(t *testing.T) {
	t.Run("Decorate both", func(t *testing.T) {
		testCase1[int]{
			Seq: func() iter.Seq[int] {
				return hiter.Decorate(slices.Values([]int{4, 4, 6, 9}), iterable.SliceAll[int]{1, 1}, iterable.SliceAll[int]{2})
			},
			Expected: []int{1, 1, 4, 2, 1, 1, 4, 2, 1, 1, 6, 2, 1, 1, 9, 2},
			BreakAt:  3,
		}.Test(t)
	})

	t.Run("Decorate nil", func(t *testing.T) {
		testCase1[int]{
			Seq: func() iter.Seq[int] {
				return hiter.Decorate(slices.Values([]int{4, 4, 6, 9}), nil, nil)
			},
			Expected: []int{4, 4, 6, 9},
			BreakAt:  3,
		}.Test(t)
	})

	t.Run("Decorate2 both", func(t *testing.T) {
		testCase2[int, int]{
			Seq: func() iter.Seq2[int, int] {
				return hiter.Decorate2(
					slices.All(
						[]int{4, 4, 6, 9}),
					hiter.KeyValues[int, int]{{2, 1}, {2, 1}},
					hiter.KeyValues[int, int]{{-1, 2}},
				)
			},
			Expected: []hiter.KeyValue[int, int]{
				{2, 1}, {2, 1}, {0, 4}, {-1, 2},
				{2, 1}, {2, 1}, {1, 4}, {-1, 2},
				{2, 1}, {2, 1}, {2, 6}, {-1, 2},
				{2, 1}, {2, 1}, {3, 9}, {-1, 2},
			},
			BreakAt: 3,
		}.Test(t)
	})

	t.Run("Decorate2 nil", func(t *testing.T) {
		testCase2[int, int]{
			Seq: func() iter.Seq2[int, int] {
				return hiter.Decorate2(slices.All([]int{4, 4, 6, 9}), nil, nil)
			},
			Expected: []hiter.KeyValue[int, int]{{0, 4}, {1, 4}, {2, 6}, {3, 9}},
			BreakAt:  3,
		}.Test(t)
	})
}
