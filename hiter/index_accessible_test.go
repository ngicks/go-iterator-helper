package hiter_test

import (
	"iter"
	"testing"

	"github.com/ngicks/go-iterator-helper/hiter"
	"github.com/ngicks/go-iterator-helper/hiter/iterable"
)

type atSliceStr []string

func (a atSliceStr) At(i int) string {
	return a[i]
}

func TestIndexAccessible(t *testing.T) {
	src := atSliceStr{
		"foo", "bar", "baz",
		"qux", "quux", "corge",
		"grault", "garply", "waldo",
		"fred", "plugh", "xyzzy",
		"thud",
	}
	indices := iterable.Range[int]{Start: 4, End: 10}
	testCase2[int, string]{
		Seq: func() iter.Seq2[int, string] {
			return hiter.IndexAccessible(src, indices.Iter())
		},
		Seqs: []func() iter.Seq2[int, string]{
			func() iter.Seq2[int, string] {
				return iterable.IndexAccessible[atSliceStr, string]{
					Atter:   src,
					Indices: indices,
				}.Iter2()
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
}
