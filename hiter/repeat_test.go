package hiter_test

import (
	"iter"
	"testing"

	"github.com/ngicks/go-iterator-helper/hiter"
	"github.com/ngicks/go-iterator-helper/hiter/iterable"
)

func TestRepeat(t *testing.T) {
	t.Run("Repeatable", func(t *testing.T) {
		testCase1[uint]{
			Seq: func() iter.Seq[uint] {
				return hiter.Repeat(uint(3), 5)
			},
			Seqs: []func() iter.Seq[uint]{
				func() iter.Seq[uint] {
					return iterable.Repeatable[uint]{V: 3, N: 5}.Iter()
				},
			},
			Expected: []uint{3, 3, 3, 3, 3},
			BreakAt:  3,
		}.Test(t)

		testCase1[uint]{
			Seq: func() iter.Seq[uint] {
				return hiter.Repeat(uint(3), -1)
			},
			Seqs: []func() iter.Seq[uint]{
				func() iter.Seq[uint] {
					return iterable.Repeatable[uint]{V: 3, N: -1}.Iter()
				},
			},
			Expected: nil,
		}.Test(t)
	})

	t.Run("Repeatable2", func(t *testing.T) {
		testCase2[uint, string]{
			Seq: func() iter.Seq2[uint, string] {
				return hiter.Repeat2(uint(7), "foo", 7)
			},
			Seqs: []func() iter.Seq2[uint, string]{
				func() iter.Seq2[uint, string] {
					return iterable.Repeatable2[uint, string]{K: 7, V: "foo", N: 7}.Iter2()
				},
			},
			Expected: []hiter.KeyValue[uint, string]{
				{7, "foo"}, {7, "foo"}, {7, "foo"},
				{7, "foo"}, {7, "foo"}, {7, "foo"},
				{7, "foo"},
			},
			BreakAt: 3,
		}.Test(t)

		testCase2[uint, string]{
			Seq: func() iter.Seq2[uint, string] {
				return hiter.Repeat2(uint(7), "foo", -1)
			},
			Seqs: []func() iter.Seq2[uint, string]{
				func() iter.Seq2[uint, string] {
					return iterable.Repeatable2[uint, string]{K: 7, V: "foo", N: -1}.Iter2()
				},
			},
		}.Test(t)
	})
}
