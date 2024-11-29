package hiter_test

import (
	"iter"
	"testing"

	"github.com/ngicks/go-iterator-helper/hiter"
	"github.com/ngicks/go-iterator-helper/hiter/internal/testcase"
	"github.com/ngicks/go-iterator-helper/hiter/iterable"
)

func TestRange(t *testing.T) {
	t.Run("uint, 0 -> 6", func(t *testing.T) {
		testcase.TestCase1[uint]{
			Seq: func() iter.Seq[uint] {
				return hiter.Range(uint(0), uint(6))
			},
			Seqs: []func() iter.Seq[uint]{
				func() iter.Seq[uint] {
					return iterable.Range[uint]{Start: 0, End: 6}.Iter()
				},
			},
			Expected: []uint{0, 1, 2, 3, 4, 5},
			BreakAt:  3,
		}.Test(t)
	})

	t.Run("float64, 0.5 -> 3.5", func(t *testing.T) {
		testcase.TestCase1[float64]{
			Seq: func() iter.Seq[float64] {
				return hiter.Range(float64(0.5), float64(3.5))
			},
			Seqs: []func() iter.Seq[float64]{
				func() iter.Seq[float64] {
					return iterable.Range[float64]{Start: 0.5, End: 3.5}.Iter()
				},
			},
			Expected: []float64{0.5, 1.5, 2.5},
			BreakAt:  1,
		}.Test(t)
	})

	t.Run("int, 4 -> -2", func(t *testing.T) {
		testcase.TestCase1[int]{
			Seq: func() iter.Seq[int] {
				return hiter.Range(4, -2)
			},
			Seqs: []func() iter.Seq[int]{
				func() iter.Seq[int] {
					return iterable.Range[int]{Start: 4, End: -2}.Iter()
				},
			},
			Expected: []int{4, 3, 2, 1, 0, -1},
			BreakAt:  3,
		}.Test(t)
	})

	t.Run("int, 1 -> 1", func(t *testing.T) {
		testcase.TestCase1[int]{
			Seq: func() iter.Seq[int] {
				return hiter.Range(1, 1)
			},
			Seqs: []func() iter.Seq[int]{
				func() iter.Seq[int] {
					return iterable.Range[int]{Start: 1, End: 1}.Iter()
				},
			},
			Expected: nil,
		}.Test(t)
	})

	t.Run("int, 5 -> 7, excludeStart, excludeEnd", func(t *testing.T) {
		testcase.TestCase1[int]{
			Seq: func() iter.Seq[int] {
				return hiter.RangeInclusive(5, 7, false, false)
			},
			Seqs: []func() iter.Seq[int]{
				func() iter.Seq[int] {
					return iterable.Range[int]{Start: 5, End: 7, ExcludesStart: true}.Iter()
				},
			},
			Expected: []int{6},
		}.Test(t)
	})

	t.Run("int, -1 -> -1, includeStart, includeEnd", func(t *testing.T) {
		testcase.TestCase1[int]{
			Seq: func() iter.Seq[int] {
				return hiter.RangeInclusive(-1, -1, true, true)
			},
			Seqs: []func() iter.Seq[int]{
				func() iter.Seq[int] {
					return iterable.Range[int]{Start: -1, End: -1, IncludesEnd: true}.Iter()
				},
			},
			Expected: []int{-1},
		}.Test(t)
	})
}
