package ih_test

import (
	"iter"
	"slices"
	"testing"

	"github.com/ngicks/go-iterator-helper/ih"
	"gotest.tools/v3/assert"
)

var (
	stringsSrc = "foobarbaz"
	runesSrc   = ".😂😎a🙂c😫"
)

func TestStringsChunk(t *testing.T) {
	t.Run("divide 9 by 3", func(t *testing.T) {
		testCase1[string]{
			Seq: func() iter.Seq[string] {
				return ih.StringsChunk(stringsSrc, 3)
			},
			Expected: []string{"foo", "bar", "baz"},
			BreakAt:  2,
		}.Test(t)
	})

	t.Run("divide 9 by 4", func(t *testing.T) {
		testCase1[string]{
			Seq: func() iter.Seq[string] {
				return ih.StringsChunk(stringsSrc, 4)
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

		testCase1[string]{
			Seq: func() iter.Seq[string] {
				return ih.StringsChunk(runesSrc, 2)
			},
			Expected: expected,
			BreakAt:  2,
		}.Test(t)

		assert.Equal(t, runesSrc, ih.StringsCollect(slices.Values(expected), len(runesSrc)))
	})

	t.Run("divide by 0", func(t *testing.T) {
		testCase1[string]{
			Seq: func() iter.Seq[string] {
				return ih.StringsChunk(runesSrc, 0)
			},
			Expected: nil,
		}.Test(t)
	})

	t.Run("divide by -1", func(t *testing.T) {
		testCase1[string]{
			Seq: func() iter.Seq[string] {
				return ih.StringsChunk(runesSrc, -1)
			},
			Expected: nil,
		}.Test(t)
	})
}

func TestStringsRuneChunk(t *testing.T) {
	t.Run("divide 9 by 3", func(t *testing.T) {
		testCase1[string]{
			Seq: func() iter.Seq[string] {
				return ih.StringsRuneChunk(stringsSrc, 3)
			},
			Expected: []string{"foo", "bar", "baz"},
			BreakAt:  2,
		}.Test(t)
	})

	t.Run("divide 9 by 4", func(t *testing.T) {
		testCase1[string]{
			Seq: func() iter.Seq[string] {
				return ih.StringsRuneChunk(stringsSrc, 4)
			},
			Expected: []string{"foob", "arba", "z"},
			BreakAt:  2,
		}.Test(t)
	})

	t.Run("divide emoji by 2", func(t *testing.T) {
		testCase1[string]{
			Seq: func() iter.Seq[string] {
				return ih.StringsRuneChunk(runesSrc, 2)
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
		testCase1[string]{
			Seq: func() iter.Seq[string] {
				return ih.StringsRuneChunk(runesSrc, 0)
			},
			Expected: nil,
		}.Test(t)
	})

	t.Run("divide emoji by -1", func(t *testing.T) {
		testCase1[string]{
			Seq: func() iter.Seq[string] {
				return ih.StringsRuneChunk(runesSrc, -1)
			},
			Expected: nil,
		}.Test(t)
	})
}

func TestStringsSplitFunc(t *testing.T) {
	const (
		src        = "foo\nbar\nbaz"
		srcCr      = "foo\n\rbar\n\rbaz\n\r"
		srcCase    = "NewHttpRequest"
		allCapital = "STOP ALL CAPITAL"
		longSingle = "foooooooooooooooooooooo"
	)

	t.Run("nil cutter", func(t *testing.T) {
		testCase1[string]{
			Seq: func() iter.Seq[string] {
				return ih.StringsSplitFunc(src, -1, nil)
			},
			Expected: []string{"foo", "bar", "baz"},
			BreakAt:  2,
		}.Test(t)

		testCase1[string]{
			Seq: func() iter.Seq[string] {
				return ih.StringsSplitFunc("", -1, nil)
			},
			Expected: nil,
		}.Test(t)
	})

	t.Run("limit by n", func(t *testing.T) {
		testCase1[string]{
			Seq: func() iter.Seq[string] {
				return ih.StringsSplitFunc(src, 1, nil)
			},
			Expected: []string{"foo", "bar\nbaz"},
			BreakAt:  1,
		}.Test(t)
	})

	t.Run("StringsCutNewLine", func(t *testing.T) {
		testCase1[string]{
			Seq: func() iter.Seq[string] {
				return ih.StringsSplitFunc(src, -1, ih.StringsCutNewLine)
			},
			Seqs: []func() iter.Seq[string]{
				func() iter.Seq[string] {
					return ih.StringsSplitFunc(srcCr, -1, ih.StringsCutNewLine)
				},
			},
			Expected: []string{"foo", "bar", "baz"},
			BreakAt:  2,
		}.Test(t)

		testCase1[string]{
			Seq: func() iter.Seq[string] {
				return ih.StringsSplitFunc(longSingle, -1, ih.StringsCutNewLine)
			},
			Expected: []string{longSingle},
		}.Test(t)
	})

	t.Run("StringsCutUpperCase", func(t *testing.T) {
		testCase1[string]{
			Seq: func() iter.Seq[string] {
				return ih.StringsSplitFunc("a", -1, ih.StringsCutUpperCase)
			},
			Expected: []string{"a"},
		}.Test(t)

		testCase1[string]{
			Seq: func() iter.Seq[string] {
				return ih.StringsSplitFunc(src, -1, ih.StringsCutUpperCase)
			},
			Expected: []string{"foo\nbar\nbaz"},
		}.Test(t)

		testCase1[string]{
			Seq: func() iter.Seq[string] {
				return ih.StringsSplitFunc(srcCase, -1, ih.StringsCutUpperCase)
			},
			Expected: []string{"New", "Http", "Request"},
		}.Test(t)

		testCase1[string]{
			Seq: func() iter.Seq[string] {
				return ih.StringsSplitFunc(allCapital, -1, ih.StringsCutUpperCase)
			},
			Expected: []string{"S", "T", "O", "P ", "A", "L", "L ", "C", "A", "P", "I", "T", "A", "L"},
		}.Test(t)
	})
}
