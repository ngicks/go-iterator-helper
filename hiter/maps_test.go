package hiter_test

import (
	"cmp"
	"iter"
	"slices"
	"testing"

	"github.com/ngicks/go-iterator-helper/hiter"
	"github.com/ngicks/go-iterator-helper/hiter/internal/testcase"
	"github.com/ngicks/go-iterator-helper/hiter/iterable"
)

func TestMaps(t *testing.T) {
	expected := map[string]string{
		"foo": "foofoo",
		"bar": "barbar",
		"baz": "bazbaz",
	}

	// Test MapsKeys function with early termination
	t.Run("MapsKeys", func(t *testing.T) {
		testcase.Two[string, string]{
			Seq: func() iter.Seq2[string, string] {
				keys := []string{"foo", "bar", "baz"}
				return hiter.MapsKeys(expected, slices.Values(keys))
			},
			Expected: []hiter.KeyValue[string, string]{
				{"foo", "foofoo"},
				{"bar", "barbar"},
				{"baz", "bazbaz"},
			},
			BreakAt: 2, // This tests the early termination case
		}.Test(t)
	})

	t.Run("MapsAll", func(t *testing.T) {
		testcase.Map[string, string]{
			Seq: func() iter.Seq2[string, string] {
				return iterable.MapAll[string, string](expected).Iter2()
			},
			Expected: expected,
			BreakAt:  2,
		}.Test(t)
	})

	t.Run("MapsAll nil", func(t *testing.T) {
		testcase.Map[string, string]{
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

	t.Run("MapsSorted", func(t *testing.T) {
		testcase.Two[string, int]{
			Seq: func() iter.Seq2[string, int] {
				return hiter.MapsSorted(expectedSI)
			},
			Seqs: []func() iter.Seq2[string, int]{
				func() iter.Seq2[string, int] {
					return iterable.MapSorted[string, int](expectedSI).Iter2()
				},
			},
			BreakAt:  2,
			Expected: []hiter.KeyValue[string, int]{{"bar", 1}, {"baz", 2}, {"foo", 0}, {"qux", 3}},
		}.Test(t)
	})

	t.Run("MapsSorted nil", func(t *testing.T) {
		testcase.Two[string, int]{
			Seq: func() iter.Seq2[string, int] {
				return hiter.MapsSorted[map[string]int](nil)
			},
			Seqs: []func() iter.Seq2[string, int]{
				func() iter.Seq2[string, int] {
					return iterable.MapSorted[string, int](nil).Iter2()
				},
			},
			Expected: nil,
		}.Test(t)
	})

	t.Run("MapsSortedFunc", func(t *testing.T) {
		testcase.Two[string, int]{
			Seq: func() iter.Seq2[string, int] {
				return hiter.MapsSortedFunc(
					expectedSI,
					func(s1, s2 string) int {
						return cmp.Compare(expectedSI[s1], expectedSI[s2])
					},
				)
			},
			Seqs: []func() iter.Seq2[string, int]{
				func() iter.Seq2[string, int] {
					return iterable.MapSortedFunc[map[string]int, string, int]{
						M: expectedSI,
						Cmp: func(s1, s2 string) int {
							return cmp.Compare(expectedSI[s1], expectedSI[s2])
						},
					}.Iter2()
				},
			},
			BreakAt: 2,
			Expected: []hiter.KeyValue[string, int]{
				{"foo", 0},
				{"bar", 1},
				{"baz", 2},
				{"qux", 3},
			},
		}.Test(t)
	})

	t.Run("MapsSortedFunc nil ", func(t *testing.T) {
		testcase.Two[string, int]{
			Seq: func() iter.Seq2[string, int] {
				return hiter.MapsSortedFunc[map[string]int](
					nil,
					func(s1, s2 string) int { return cmp.Compare(expectedSI[s1], expectedSI[s2]) },
				)
			},
			Seqs: []func() iter.Seq2[string, int]{
				func() iter.Seq2[string, int] {
					return iterable.MapSortedFunc[map[string]int, string, int]{
						M:   nil,
						Cmp: func(s1, s2 string) int { return cmp.Compare(expectedSI[s1], expectedSI[s2]) },
					}.Iter2()
				},
			},
			Expected: nil,
		}.Test(t)
	})

	t.Run("MapsOverlay", func(t *testing.T) {
		m1 := map[string]int{
			"foo":   1,
			"bar":   2,
			"baz":   3,
			"qux":   4,
			"quux":  5,
			"corge": 6,
		}
		m2 := map[string]int{
			"foo": 7,
			"bar": 8,
			"baz": 9,
			"qux": 10,

			"whooooooaa": 0,
		}
		m3 := map[string]int{
			"foo": 11,
			"bar": 12,

			"weeeeeeeee": 0,
		}
		testcase.Map[string, int]{
			Seq: func() iter.Seq2[string, int] {
				return hiter.MapsOverlay(m1, m2, m3)
			},
			Seqs: []func() iter.Seq2[string, int]{
				func() iter.Seq2[string, int] {
					return iterable.MapOverlay[map[string]int, string, int]{
						Maps: []map[string]int{m1, m2, m3},
					}.Iter2()
				},
			},
			Expected: map[string]int{
				"foo":        11,
				"bar":        12,
				"baz":        9,
				"qux":        10,
				"quux":       5,
				"corge":      6,
				"whooooooaa": 0,
				"weeeeeeeee": 0,
			},
		}.Test(t)
	})
}
