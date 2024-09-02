package ih_test

import (
	"cmp"
	"iter"
	"slices"
	"testing"

	"github.com/ngicks/go-iterator-helper/ih"
	"github.com/ngicks/go-iterator-helper/ih/iterable"
)

func TestMap(t *testing.T) {
	expected := map[string]string{
		"foo": "foofoo",
		"bar": "barbar",
		"baz": "bazbaz",
	}

	t.Run("MapAll", func(t *testing.T) {
		testCaseMap[string, string]{
			Seq: func() iter.Seq2[string, string] {
				return iterable.MapAll[string, string](expected).Iter2()
			},
			Expected: expected,
			BreakAt:  2,
		}.Test(t)
	})

	t.Run("MapAll nil", func(t *testing.T) {
		testCaseMap[string, string]{
			Seq: func() iter.Seq2[string, string] {
				return iterable.MapAll[string, string](nil).Iter2()
			},
			Expected: map[string]string{},
		}.Test(t)
	})

	expectedSI := map[string]int{
		"foo": 0,
		"bar": 1,
		"baz": 2,
		"qux": 3,
	}

	t.Run("MapSorted", func(t *testing.T) {
		testCase2[string, int]{
			Seq: func() iter.Seq2[string, int] {
				return iterable.MapSorted[string, int](expectedSI).Iter2()
			},
			BreakAt:  2,
			Expected: []ih.KeyValue[string, int]{{"bar", 1}, {"baz", 2}, {"foo", 0}, {"qux", 3}},
		}.Test(t)
	})

	t.Run("MapSorted nil", func(t *testing.T) {
		testCase2[string, int]{
			Seq: func() iter.Seq2[string, int] {
				return iterable.MapSorted[string, int](nil).Iter2()
			},
			Expected: nil,
		}.Test(t)
	})

	t.Run("MapSortedFunc", func(t *testing.T) {
		testCase2[string, int]{
			Seq: func() iter.Seq2[string, int] {
				return iterable.MapSortedFunc[map[string]int, string, int]{
					M:   expectedSI,
					Cmp: func(s1, s2 string) int { return cmp.Compare(expectedSI[s1], expectedSI[s2]) },
				}.Iter2()
			},
			BreakAt: 2,
			Expected: []ih.KeyValue[string, int]{
				{"foo", 0},
				{"bar", 1},
				{"baz", 2},
				{"qux", 3},
			},
		}.Test(t)
	})

	t.Run("MapSortedFunc nil ", func(t *testing.T) {
		testCase2[string, int]{
			Seq: func() iter.Seq2[string, int] {
				return iterable.MapSortedFunc[map[string]int, string, int]{
					M:   nil,
					Cmp: func(s1, s2 string) int { return cmp.Compare(expectedSI[s1], expectedSI[s2]) },
				}.Iter2()
			},
			Expected: nil,
		}.Test(t)
	})
}

func TestSlice(t *testing.T) {
	t.Run("SliceAll", func(t *testing.T) {
		testCase1[int]{
			Seq: func() iter.Seq[int] {
				return iterable.SliceAll[int](srcInt1).Iter()
			},
			BreakAt:  2,
			Expected: srcInt1,
		}.Test(t)
		testCase2[int, int]{
			Seq: func() iter.Seq2[int, int] {
				return iterable.SliceAll[int](srcInt1).Iter2()
			},
			BreakAt:  2,
			Expected: []ih.KeyValue[int, int]{{0, 12}, {1, 76}, {2, 8}, {3, 9}, {4, 923}},
		}.Test(t)
	})

	t.Run("SliceBackward", func(t *testing.T) {
		testCase1[int]{
			Seq: func() iter.Seq[int] {
				return iterable.SliceBackward[int](srcInt1).Iter()
			},
			BreakAt:  2,
			Expected: slices.Collect(ih.OmitF(slices.Backward(srcInt1))),
		}.Test(t)
		testCase2[int, int]{
			Seq: func() iter.Seq2[int, int] {
				return iterable.SliceBackward[int](srcInt1).Iter2()
			},
			BreakAt:  2,
			Expected: ih.Collect2(slices.Backward(srcInt1)),
		}.Test(t)
	})
}
