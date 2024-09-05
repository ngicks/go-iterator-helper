package hiter_test

import (
	"iter"
	"testing"

	"github.com/ngicks/go-iterator-helper/hiter"
	"github.com/ngicks/go-iterator-helper/hiter/iterable"
)

func TestRepeat(t *testing.T) {
	// avoiding depending to xiter.
	limit10 := func(seq iter.Seq[uint]) iter.Seq[uint] {
		var count int
		return hiter.LimitUntil(seq, func(uint) bool { count++; return count <= 10 })
	}
	limit102 := func(seq iter.Seq2[uint, string]) iter.Seq2[uint, string] {
		var count int
		return hiter.LimitUntil2(seq, func(uint, string) bool { count++; return count <= 10 })
	}

	t.Run("Repeatable", func(t *testing.T) {
		testCase1[uint]{
			Seq: func() iter.Seq[uint] {
				return hiter.Repeat(uint(3), 5)
			},
			Seqs: []func() iter.Seq[uint]{
				func() iter.Seq[uint] {
					return iterable.Repeatable[uint]{V: 3, N: 5}.Iter()
				},
				func() iter.Seq[uint] {
					return hiter.RepeatFunc(func() uint { return 3 }, 5)
				},
				func() iter.Seq[uint] {
					return iterable.RepeatableFunc[uint]{FnV: func() uint { return 3 }, N: 5}.Iter()
				},
			},
			Expected: []uint{3, 3, 3, 3, 3},
			BreakAt:  3,
		}.Test(t)

		testCase1[uint]{
			Seq: func() iter.Seq[uint] {
				return limit10(hiter.Repeat(uint(3), -1))
			},
			Seqs: []func() iter.Seq[uint]{
				func() iter.Seq[uint] {
					return limit10(iterable.Repeatable[uint]{V: 3, N: -1}.Iter())
				},
				func() iter.Seq[uint] {
					return limit10(hiter.RepeatFunc(func() uint { return 3 }, -1))
				},
				func() iter.Seq[uint] {
					return limit10(iterable.RepeatableFunc[uint]{FnV: func() uint { return 3 }, N: -1}.Iter())
				},
			},
			Expected: []uint{
				3, 3, 3, 3, 3,
				3, 3, 3, 3, 3,
			},
			BreakAt: 5,
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
				func() iter.Seq2[uint, string] {
					return hiter.RepeatFunc2(func() uint { return 7 }, func() string { return "foo" }, 7)
				},
				func() iter.Seq2[uint, string] {
					return iterable.RepeatableFunc2[uint, string]{FnK: func() uint { return 7 }, FnV: func() string { return "foo" }, N: 7}.Iter2()
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
				return limit102(hiter.Repeat2(uint(7), "foo", -1))
			},
			Seqs: []func() iter.Seq2[uint, string]{
				func() iter.Seq2[uint, string] {
					return limit102(iterable.Repeatable2[uint, string]{K: 7, V: "foo", N: -1}.Iter2())
				},
				func() iter.Seq2[uint, string] {
					return limit102(hiter.RepeatFunc2(func() uint { return 7 }, func() string { return "foo" }, -1))
				},
				func() iter.Seq2[uint, string] {
					return limit102(iterable.RepeatableFunc2[uint, string]{FnK: func() uint { return 7 }, FnV: func() string { return "foo" }, N: -1}.Iter2())
				},
			},
			Expected: hiter.KeyValues[uint, string]{
				{7, "foo"}, {7, "foo"}, {7, "foo"}, {7, "foo"}, {7, "foo"},
				{7, "foo"}, {7, "foo"}, {7, "foo"}, {7, "foo"}, {7, "foo"},
			},
		}.Test(t)
	})
}
