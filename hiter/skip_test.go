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
				return hiter.Skip(5, slices.Values(src))
			},
			Expected: []int{1, 4, 2, 1, 1, 6, 2, 1, 1, 9, 2},
			BreakAt:  3,
		}.Test(t)
		testCase1[int]{
			Seq: func() iter.Seq[int] {
				return hiter.Skip(5, slices.Values(src[:2]))
			},
		}.Test(t)
	})

	t.Run("SkipLast", func(t *testing.T) {
		testCase1[int]{
			Seq: func() iter.Seq[int] {
				return hiter.SkipLast(5, slices.Values(src))
			},
			Expected: []int{1, 1, 4, 2, 1, 1, 4, 2, 1, 1, 6},
			BreakAt:  3,
		}.Test(t)
		testCase1[int]{
			Seq: func() iter.Seq[int] {
				return hiter.SkipLast(5, slices.Values(src[:2]))
			},
		}.Test(t)
	})

	t.Run("Skip2", func(t *testing.T) {
		testCase2[int, int]{
			Seq: func() iter.Seq2[int, int] {
				return hiter.Skip2(5, slices.All(src))
			},
			Expected: []hiter.KeyValue[int, int]{
				{5, 1}, {6, 4}, {7, 2},
				{8, 1}, {9, 1}, {10, 6},
				{11, 2}, {12, 1}, {13, 1},
				{14, 9}, {15, 2},
			},
			BreakAt: 3,
		}.Test(t)
		testCase2[int, int]{
			Seq: func() iter.Seq2[int, int] {
				return hiter.Skip2(5, slices.All(src[:2]))
			},
		}.Test(t)
	})

	t.Run("SkipLast2", func(t *testing.T) {
		testCase2[int, int]{
			Seq: func() iter.Seq2[int, int] {
				return hiter.SkipLast2(5, slices.All(src))
			},
			Expected: []hiter.KeyValue[int, int]{{0, 1}, {1, 1}, {2, 4}, {3, 2}, {4, 1}, {5, 1}, {6, 4}, {7, 2}, {8, 1}, {9, 1}, {10, 6}},
			BreakAt:  3,
		}.Test(t)
		testCase2[int, int]{
			Seq: func() iter.Seq2[int, int] {
				return hiter.SkipLast2(5, slices.All(src[:2]))
			},
		}.Test(t)
	})
}

func TestSkipUntil(t *testing.T) {
	src := []int{1, 1, 4, 2, 1, 1, 4, 2, 1, 1, 6, 2, 1, 1, 9, 2}
	t.Run("SkipWhile", func(t *testing.T) {
		testCase1[int]{
			Seq: func() iter.Seq[int] {
				return hiter.SkipWhile(func(e int) bool { return e == 6 }, slices.Values(src))
			},
			Expected: []int{6, 2, 1, 1, 9, 2},
			BreakAt:  3,
		}.Test(t)
	})

	t.Run("SkipWhile2", func(t *testing.T) {
		testCase2[int, int]{
			Seq: func() iter.Seq2[int, int] {
				return hiter.SkipWhile2(func(i, e int) bool { return e == 6 }, slices.All(src))
			},
			Expected: []hiter.KeyValue[int, int]{
				{10, 6}, {11, 2}, {12, 1},
				{13, 1}, {14, 9}, {15, 2},
			},
			BreakAt: 3,
		}.Test(t)
	})
}
