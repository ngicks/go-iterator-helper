package hiter_test

import (
	"iter"
	"slices"
	"testing"

	"github.com/ngicks/go-iterator-helper/hiter"
	"github.com/ngicks/go-iterator-helper/hiter/internal/testcase"
)

func TestFlatten(t *testing.T) {
	var (
		nested1 = slices.Values([][]int{{0, 1, 2}, {3, 4}, {5}})
		nested2 = slices.Values([][]int{{10}, {7, 8}, {1}})
		flat1   = slices.Values([]int{10, 7, 8, 1})
		flat2   = slices.Values([]int{0, 1, 2, 3, 4, 5})
	)

	flattenResult := []int{0, 1, 2, 3, 4, 5}
	t.Run("Flatten", func(t *testing.T) {
		testcase.One[int]{
			Seq: func() iter.Seq[int] {
				return hiter.Flatten(nested1)
			},
			Expected: flattenResult,
			BreakAt:  4,
		}.Test(t)
	})
	t.Run("Flatten", func(t *testing.T) {
		testcase.One[int]{
			Seq: func() iter.Seq[int] {
				return hiter.FlattenSeq(hiter.Map(slices.Values, nested1))
			},
			Expected: flattenResult,
			BreakAt:  4,
		}.Test(t)
	})

	t.Run("FlattenSeq2", func(t *testing.T) {
		testcase.Two[int, int]{
			Seq: func() iter.Seq2[int, int] {
				return hiter.FlattenSeq2(slices.Values(
					[]iter.Seq2[int, int]{
						hiter.Values2([]hiter.KeyValue[int, int]{{0, 1}, {2, 2}}),
						hiter.Values2([]hiter.KeyValue[int, int]{{1, 0}, {4, 5}, {2, 9}}),
					},
				))
			},
			Expected: []hiter.KeyValue[int, int]{{0, 1}, {2, 2}, {1, 0}, {4, 5}, {2, 9}},
			BreakAt:  2,
		}.Test(t)
	})

	flattenFResult := []hiter.KeyValue[int, int]{{0, 10}, {1, 10}, {2, 10}, {3, 7}, {4, 7}, {5, 8}}
	t.Run("FlattenF", func(t *testing.T) {
		testcase.Two[int, int]{
			Seq: func() iter.Seq2[int, int] {
				return hiter.FlattenF(hiter.Pairs(
					nested1,
					flat1,
				))
			},
			Expected: flattenFResult,
			BreakAt:  4,
		}.Test(t)
	})
	t.Run("FlattenSeqF", func(t *testing.T) {
		testcase.Two[int, int]{
			Seq: func() iter.Seq2[int, int] {
				return hiter.FlattenSeqF(hiter.Pairs(
					hiter.Map(slices.Values, nested1),
					flat1,
				))
			},
			Expected: flattenFResult,
			BreakAt:  4,
		}.Test(t)
	})

	flattenLResult := []hiter.KeyValue[int, int]{{0, 10}, {1, 7}, {1, 8}, {2, 1}}
	t.Run("FlattenL", func(t *testing.T) {
		testcase.Two[int, int]{
			Seq: func() iter.Seq2[int, int] {
				return hiter.FlattenL(hiter.Pairs(
					flat2,
					nested2,
				))
			},
			Expected: flattenLResult,
			BreakAt:  2,
		}.Test(t)
	})
	t.Run("FlattenSeqL", func(t *testing.T) {
		testcase.Two[int, int]{
			Seq: func() iter.Seq2[int, int] {
				return hiter.FlattenSeqL(hiter.Pairs(
					flat2,
					hiter.Map(slices.Values, nested2),
				))
			},
			Expected: flattenLResult,
			BreakAt:  2,
		}.Test(t)
	})
}
