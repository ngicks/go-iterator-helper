package hiter_test

import (
	"iter"
	"slices"
	"testing"

	"github.com/ngicks/go-iterator-helper/hiter"
)

func TestAlternate(t *testing.T) {
	t.Run("Alternate", func(t *testing.T) {
		t.Run("0", func(t *testing.T) {
			testCase1[int]{
				Seq: func() iter.Seq[int] {
					return hiter.Alternate[int]()
				},
			}.Test(t)
		})

		t.Run("only 1", func(t *testing.T) {
			testCase1[int]{
				Seq: func() iter.Seq[int] {
					return hiter.Alternate(
						slices.Values([]int{1, 2, 3}),
					)
				},
				Expected: []int{1, 2, 3},
				BreakAt:  2,
			}.Test(t)
		})

		t.Run("2", func(t *testing.T) {
			testCase1[int]{
				Seq: func() iter.Seq[int] {
					return hiter.Alternate(
						slices.Values([]int{1, 2, 3}),
						slices.Values([]int{4, 5, 6}),
					)
				},
				Expected: []int{1, 4, 2, 5, 3, 6},
				BreakAt:  2,
			}.Test(t)
		})

		t.Run("3", func(t *testing.T) {
			testCase1[int]{
				Seq: func() iter.Seq[int] {
					return hiter.Alternate(
						slices.Values([]int{1, 2, 3}),
						slices.Values([]int{4, 5, 6}),
						slices.Values([]int{7, 8, 9}),
					)
				},
				Expected: []int{1, 4, 7, 2, 5, 8, 3, 6, 9},
				BreakAt:  5,
			}.Test(t)
		})

		t.Run("first shorter", func(t *testing.T) {
			testCase1[int]{
				Seq: func() iter.Seq[int] {
					return hiter.Alternate(
						slices.Values([]int{1}),
						slices.Values([]int{4, 5, 6}),
					)
				},
				Expected: []int{1, 4},
				BreakAt:  2,
			}.Test(t)
		})

		t.Run("second shorter", func(t *testing.T) {
			testCase1[int]{
				Seq: func() iter.Seq[int] {
					return hiter.Alternate(
						slices.Values([]int{1, 2, 3}),
						slices.Values([]int{4, 5}),
					)
				},
				Expected: []int{1, 4, 2, 5, 3},
				BreakAt:  2,
			}.Test(t)
		})
	})

	t.Run("Alternate2", func(t *testing.T) {
		t.Run("0", func(t *testing.T) {
			testCase2[int, int]{
				Seq: func() iter.Seq2[int, int] {
					return hiter.Alternate2[int, int]()
				},
			}.Test(t)
		})

		t.Run("only 1", func(t *testing.T) {
			testCase2[int, int]{
				Seq: func() iter.Seq2[int, int] {
					return hiter.Alternate2(
						slices.All([]int{1, 2, 3}),
					)
				},
				Expected: []hiter.KeyValue[int, int]{{0, 1}, {1, 2}, {2, 3}},
				BreakAt:  2,
			}.Test(t)
		})

		t.Run("2", func(t *testing.T) {
			testCase2[int, int]{
				Seq: func() iter.Seq2[int, int] {
					return hiter.Alternate2(
						slices.All([]int{1, 2, 3}),
						slices.All([]int{4, 5, 6}),
					)
				},
				Expected: []hiter.KeyValue[int, int]{{0, 1}, {0, 4}, {1, 2}, {1, 5}, {2, 3}, {2, 6}},
				BreakAt:  2,
			}.Test(t)
		})

		t.Run("3", func(t *testing.T) {
			testCase2[int, int]{
				Seq: func() iter.Seq2[int, int] {
					return hiter.Alternate2(
						slices.All([]int{1, 2, 3}),
						slices.All([]int{4, 5, 6}),
						slices.All([]int{7, 8, 9}),
					)
				},
				Expected: []hiter.KeyValue[int, int]{
					{0, 1}, {0, 4}, {0, 7},
					{1, 2}, {1, 5}, {1, 8},
					{2, 3}, {2, 6}, {2, 9},
				},
				BreakAt: 5,
			}.Test(t)
		})

		t.Run("first shorter", func(t *testing.T) {
			testCase2[int, int]{
				Seq: func() iter.Seq2[int, int] {
					return hiter.Alternate2(
						slices.All([]int{1}),
						slices.All([]int{4, 5, 6}),
					)
				},
				Expected: []hiter.KeyValue[int, int]{{0, 1}, {0, 4}},
				BreakAt:  2,
			}.Test(t)
		})

		t.Run("second shorter", func(t *testing.T) {
			testCase2[int, int]{
				Seq: func() iter.Seq2[int, int] {
					return hiter.Alternate2(
						slices.All([]int{1, 2, 3}),
						slices.All([]int{4, 5}),
					)
				},
				Expected: []hiter.KeyValue[int, int]{
					{0, 1}, {0, 4},
					{1, 2}, {1, 5},
					{2, 3},
				},
				BreakAt: 2,
			}.Test(t)
		})
	})
}
