package hiter_test

import (
	"iter"
	"slices"
	"testing"

	"github.com/ngicks/go-iterator-helper/hiter"
)

func TestSkip(t *testing.T) {
	src := []int{1, 1, 4, 2, 1, 1, 4, 2, 1, 1, 6, 2, 1, 1, 9, 2}
	t.Run("Skip", func(t *testing.T) {
		testCase1[int]{
			Seq: func() iter.Seq[int] {
				return hiter.Skip(slices.Values(src), 5)
			},
			Expected: []int{1, 4, 2, 1, 1, 6, 2, 1, 1, 9, 2},
			BreakAt:  3,
		}.Test(t)
	})

	t.Run("Skip2", func(t *testing.T) {
		testCase2[int, int]{
			Seq: func() iter.Seq2[int, int] {
				return hiter.Skip2(slices.All(src), 5)
			},
			Expected: []hiter.KeyValue[int, int]{
				{5, 1}, {6, 4}, {7, 2},
				{8, 1}, {9, 1}, {10, 6},
				{11, 2}, {12, 1}, {13, 1},
				{14, 9}, {15, 2},
			},
			BreakAt: 3,
		}.Test(t)
	})
}

func TestSkipUntil(t *testing.T) {
	src := []int{1, 1, 4, 2, 1, 1, 4, 2, 1, 1, 6, 2, 1, 1, 9, 2}
	t.Run("SkipWhile", func(t *testing.T) {
		testCase1[int]{
			Seq: func() iter.Seq[int] {
				return hiter.SkipWhile(slices.Values(src), func(e int) bool { return e == 6 })
			},
			Expected: []int{6, 2, 1, 1, 9, 2},
			BreakAt:  3,
		}.Test(t)
	})

	t.Run("SkipWhile2", func(t *testing.T) {
		testCase2[int, int]{
			Seq: func() iter.Seq2[int, int] {
				return hiter.SkipWhile2(slices.All(src), func(i, e int) bool { return e == 6 })
			},
			Expected: []hiter.KeyValue[int, int]{
				{10, 6}, {11, 2}, {12, 1},
				{13, 1}, {14, 9}, {15, 2},
			},
			BreakAt: 3,
		}.Test(t)
	})
}
