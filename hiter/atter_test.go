package hiter_test

import (
	"iter"
	"slices"
	"testing"

	"github.com/ngicks/go-iterator-helper/hiter"
	"github.com/ngicks/go-iterator-helper/hiter/iterable"
)

func TestAtter(t *testing.T) {
	indices := iterable.Range[int]{Start: 4, End: 10}

	testCase2[int, string]{
		Seq: func() iter.Seq2[int, string] {
			return hiter.AtterIndices(atSliceSrc, slices.Values([]int{2, 6, 1, 2}))
		},
		Expected: []hiter.KeyValue[int, string]{
			{2, "baz"},
			{6, "grault"},
			{1, "bar"},
			{2, "baz"},
		},
		BreakAt: 3,
	}.Test(t)

	testCase2[int, string]{
		Seq: func() iter.Seq2[int, string] {
			return iterable.Atter[atSliceStr, string]{
				Atter:   atSliceSrc,
				Indices: indices,
			}.Iter2()
		},
		Seqs: []func() iter.Seq2[int, string]{
			func() iter.Seq2[int, string] {
				return hiter.AtterRange(atSliceSrc, 4, 10)
			},
		},
		Expected: []hiter.KeyValue[int, string]{
			{4, "quux"},
			{5, "corge"},
			{6, "grault"},
			{7, "garply"},
			{8, "waldo"},
			{9, "fred"},
		},
		BreakAt: 3,
	}.Test(t)

	testCase2[int, string]{
		Seq: func() iter.Seq2[int, string] {
			return hiter.AtterAll(atSliceSrc)
		},
		Expected: []hiter.KeyValue[int, string]{
			{0, "foo"},
			{1, "bar"},
			{2, "baz"},
			{3, "qux"},
			{4, "quux"},
			{5, "corge"},
			{6, "grault"},
			{7, "garply"},
			{8, "waldo"},
			{9, "fred"},
			{10, "plugh"},
			{11, "xyzzy"},
			{12, "thud"},
		},
		BreakAt: 3,
	}.Test(t)
}
