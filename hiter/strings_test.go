package hiter_test

import (
	"iter"
	"slices"
	"testing"

	"github.com/ngicks/go-iterator-helper/hiter"
	"github.com/ngicks/go-iterator-helper/hiter/internal/testcase"
	"gotest.tools/v3/assert"
)

var (
	stringsSrc = "foobarbaz"
	runesSrc   = ".😂😎a🙂c😫"
)

func TestStringsChunk(t *testing.T) {
	t.Run("divide 9 by 3", func(t *testing.T) {
		testcase.One[string]{
			Seq: func() iter.Seq[string] {
				return hiter.StringsChunk(stringsSrc, 3)
			},
			Expected: []string{"foo", "bar", "baz"},
			BreakAt:  2,
		}.Test(t)
	})

	t.Run("divide 9 by 4", func(t *testing.T) {
		testcase.One[string]{
			Seq: func() iter.Seq[string] {
				return hiter.StringsChunk(stringsSrc, 4)
			},
			Expected: []string{"foob", "arba", "z"},
			BreakAt:  2,
		}.Test(t)
	})

	t.Run("divide emoji by 2", func(t *testing.T) {
		expected := []string{
			".\xf0",
			"\x9f\x98",
			"\x82\xf0",
			"\x9f\x98",
			"\x8ea",
			"\xf0\x9f",
			"\x99\x82",
			"c\xf0",
			"\x9f\x98",
			"\xab",
		}

		testcase.One[string]{
			Seq: func() iter.Seq[string] {
				return hiter.StringsChunk(runesSrc, 2)
			},
			Expected: expected,
			BreakAt:  2,
		}.Test(t)

		assert.Equal(t, runesSrc, hiter.StringsCollect(len(runesSrc), slices.Values(expected)))
	})

	t.Run("divide by 0", func(t *testing.T) {
		testcase.One[string]{
			Seq: func() iter.Seq[string] {
				return hiter.StringsChunk(runesSrc, 0)
			},
			Expected: nil,
		}.Test(t)
	})

	t.Run("divide by -1", func(t *testing.T) {
		testcase.One[string]{
			Seq: func() iter.Seq[string] {
				return hiter.StringsChunk(runesSrc, -1)
			},
			Expected: nil,
		}.Test(t)
	})
}

func TestStringsRuneChunk(t *testing.T) {
	t.Run("divide 9 by 3", func(t *testing.T) {
		testcase.One[string]{
			Seq: func() iter.Seq[string] {
				return hiter.StringsRuneChunk(stringsSrc, 3)
			},
			Expected: []string{"foo", "bar", "baz"},
			BreakAt:  2,
		}.Test(t)
	})

	t.Run("divide 9 by 4", func(t *testing.T) {
		testcase.One[string]{
			Seq: func() iter.Seq[string] {
				return hiter.StringsRuneChunk(stringsSrc, 4)
			},
			Expected: []string{"foob", "arba", "z"},
			BreakAt:  2,
		}.Test(t)
	})

	t.Run("divide emoji by 2", func(t *testing.T) {
		testcase.One[string]{
			Seq: func() iter.Seq[string] {
				return hiter.StringsRuneChunk(runesSrc, 2)
			},
			Expected: []string{
				".😂",
				"😎a",
				"🙂c",
				"😫",
			},
			BreakAt: 2,
		}.Test(t)
	})

	t.Run("divide emoji by 0", func(t *testing.T) {
		testcase.One[string]{
			Seq: func() iter.Seq[string] {
				return hiter.StringsRuneChunk(runesSrc, 0)
			},
			Expected: nil,
		}.Test(t)
	})

	t.Run("divide emoji by -1", func(t *testing.T) {
		testcase.One[string]{
			Seq: func() iter.Seq[string] {
				return hiter.StringsRuneChunk(runesSrc, -1)
			},
			Expected: nil,
		}.Test(t)
	})
}

func TestStringsSplitFunc(t *testing.T) {
	const (
		src        = "foo\nbar\nbaz"
		srcCr      = "foo\r\nbar\r\nbaz\r\n"
		srcCase    = "NewHttpRequest"
		allCapital = "STOP ALL CAPITAL"
		longSingle = "foooooooooooooooooooooo"
	)

	t.Run("nil cutter", func(t *testing.T) {
		testcase.One[string]{
			Seq: func() iter.Seq[string] {
				return hiter.StringsSplitFunc(src, -1, nil)
			},
			Expected: []string{"foo", "bar", "baz"},
			BreakAt:  2,
		}.Test(t)

		testcase.One[string]{
			Seq: func() iter.Seq[string] {
				return hiter.StringsSplitFunc("", -1, nil)
			},
			Expected: nil,
		}.Test(t)
	})

	t.Run("limit by n", func(t *testing.T) {
		testcase.One[string]{
			Seq: func() iter.Seq[string] {
				return hiter.StringsSplitFunc(src, 1, nil)
			},
			Expected: []string{"foo", "bar\nbaz"},
			BreakAt:  1,
		}.Test(t)
	})

	t.Run("StringsCutNewLine", func(t *testing.T) {
		testcase.One[string]{
			Seq: func() iter.Seq[string] {
				return hiter.StringsSplitFunc(src, -1, hiter.StringsCutNewLine)
			},
			Seqs: []func() iter.Seq[string]{
				func() iter.Seq[string] {
					return hiter.StringsSplitFunc(srcCr, -1, hiter.StringsCutNewLine)
				},
			},
			Expected: []string{"foo", "bar", "baz"},
			BreakAt:  2,
		}.Test(t)

		testcase.One[string]{
			Seq: func() iter.Seq[string] {
				return hiter.StringsSplitFunc(longSingle, -1, hiter.StringsCutNewLine)
			},
			Expected: []string{longSingle},
		}.Test(t)
	})

	t.Run("StringsCutUpperCase", func(t *testing.T) {
		testcase.One[string]{
			Seq: func() iter.Seq[string] {
				return hiter.StringsSplitFunc("a", -1, hiter.StringsCutUpperCase)
			},
			Expected: []string{"a"},
		}.Test(t)

		testcase.One[string]{
			Seq: func() iter.Seq[string] {
				return hiter.StringsSplitFunc(src, -1, hiter.StringsCutUpperCase)
			},
			Expected: []string{"foo\nbar\nbaz"},
		}.Test(t)

		testcase.One[string]{
			Seq: func() iter.Seq[string] {
				return hiter.StringsSplitFunc(srcCase, -1, hiter.StringsCutUpperCase)
			},
			Expected: []string{"New", "Http", "Request"},
		}.Test(t)

		testcase.One[string]{
			Seq: func() iter.Seq[string] {
				return hiter.StringsSplitFunc(allCapital, -1, hiter.StringsCutUpperCase)
			},
			Expected: []string{"S", "T", "O", "P ", "A", "L", "L ", "C", "A", "P", "I", "T", "A", "L"},
		}.Test(t)
	})

	t.Run("StringsCutWord", func(t *testing.T) {
		testcase.One[string]{
			Seq: func() iter.Seq[string] {
				return hiter.StringsSplitFunc(allCapital, -1, hiter.StringsCutWord)
			},
			Expected: []string{"STOP", "ALL", "CAPITAL"},
		}.Test(t)
		testcase.One[string]{
			Seq: func() iter.Seq[string] {
				return hiter.StringsSplitFunc("AAA\t\tBBB    CCC    ", -1, hiter.StringsCutWord)
			},
			Expected: []string{"AAA", "BBB", "CCC"},
		}.Test(t)
	})
}
