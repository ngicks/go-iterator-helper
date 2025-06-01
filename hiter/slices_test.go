package hiter_test

import (
	"iter"
	"slices"
	"testing"

	"github.com/ngicks/go-iterator-helper/hiter"
	"github.com/ngicks/go-iterator-helper/hiter/internal/testcase"
	"github.com/ngicks/go-iterator-helper/hiter/iterable"
)

func TestSlice(t *testing.T) {
	t.Run("SliceAll", func(t *testing.T) {
		testcase.One[int]{
			Seq: func() iter.Seq[int] {
				return iterable.SliceAll[int](srcInt1).Iter()
			},
			BreakAt:  2,
			Expected: srcInt1,
		}.Test(t)
		testcase.Two[int, int]{
			Seq: func() iter.Seq2[int, int] {
				return iterable.SliceAll[int](srcInt1).Iter2()
			},
			BreakAt:  2,
			Expected: []hiter.KeyValue[int, int]{{0, 12}, {1, 76}, {2, 8}, {3, 9}, {4, 923}},
		}.Test(t)
	})

	t.Run("SliceBackward", func(t *testing.T) {
		testcase.One[int]{
			Seq: func() iter.Seq[int] {
				return iterable.SliceBackward[int](srcInt1).Iter()
			},
			BreakAt:  2,
			Expected: slices.Collect(hiter.OmitF(slices.Backward(srcInt1))),
		}.Test(t)
		testcase.Two[int, int]{
			Seq: func() iter.Seq2[int, int] {
				return iterable.SliceBackward[int](srcInt1).Iter2()
			},
			BreakAt:  2,
			Expected: hiter.Collect2(slices.Backward(srcInt1)),
		}.Test(t)
	})
}
