package hiter_test

import (
	"iter"
	"slices"
	"testing"

	"github.com/ngicks/go-iterator-helper/hiter"
	"github.com/ngicks/go-iterator-helper/hiter/internal/testcase"
)

var (
	limitUntilSrc = []int{1, 1, 4, 2, 1, 1, 4, 2, 1, 1, 6, 2, 1, 1, 9, 2}
)

func TestLimitUntil(t *testing.T) {
	testcase.TestCase1[int]{
		Seq: func() iter.Seq[int] {
			return hiter.LimitUntil(func(e int) bool { return e < 6 }, slices.Values(limitUntilSrc))
		},
		Expected: []int{1, 1, 4, 2, 1, 1, 4, 2, 1, 1},
		BreakAt:  3,
	}.Test(t)
}

func TestLimitAfter(t *testing.T) {
	testcase.TestCase1[int]{
		Seq: func() iter.Seq[int] {
			return hiter.LimitAfter(func(e int) bool { return e < 6 }, slices.Values(limitUntilSrc))
		},
		Expected: []int{1, 1, 4, 2, 1, 1, 4, 2, 1, 1, 6},
		BreakAt:  3,
	}.Test(t)
}

func TestLimitUntil2(t *testing.T) {

	t.Run("limit by value", func(t *testing.T) {
		testcase.TestCase2[int, int]{
			Seq: func() iter.Seq2[int, int] {
				return hiter.LimitUntil2(func(i, e int) bool { return e < 6 }, slices.All(limitUntilSrc))
			},
			Expected: []hiter.KeyValue[int, int]{{0, 1}, {1, 1}, {2, 4}, {3, 2}, {4, 1}, {5, 1}, {6, 4}, {7, 2}, {8, 1}, {9, 1}},
			BreakAt:  3,
		}.Test(t)
	})

	t.Run("limit by key", func(t *testing.T) {
		testcase.TestCase2[int, int]{
			Seq: func() iter.Seq2[int, int] {
				return hiter.LimitUntil2(func(i, e int) bool { return i < 3 }, slices.All(limitUntilSrc))
			},
			Expected: []hiter.KeyValue[int, int]{{0, 1}, {1, 1}, {2, 4}},
			BreakAt:  2,
		}.Test(t)
	})
}

func TestLimitAfter2(t *testing.T) {

	t.Run("limit by value", func(t *testing.T) {
		testcase.TestCase2[int, int]{
			Seq: func() iter.Seq2[int, int] {
				return hiter.LimitAfter2(func(i, e int) bool { return e < 4 }, slices.All(limitUntilSrc))
			},
			Expected: []hiter.KeyValue[int, int]{{0, 1}, {1, 1}, {2, 4}},
			BreakAt:  3,
		}.Test(t)
	})

	t.Run("limit by key", func(t *testing.T) {
		testcase.TestCase2[int, int]{
			Seq: func() iter.Seq2[int, int] {
				return hiter.LimitAfter2(func(i, e int) bool { return i < 3 }, slices.All(limitUntilSrc))
			},
			Expected: []hiter.KeyValue[int, int]{{0, 1}, {1, 1}, {2, 4}, {3, 2}},
			BreakAt:  2,
		}.Test(t)
	})
}
