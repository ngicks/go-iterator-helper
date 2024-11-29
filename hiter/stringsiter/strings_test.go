package stringsiter_test

import (
	"iter"
	"slices"
	"testing"

	"github.com/ngicks/go-iterator-helper/hiter/internal/testcase"
	"github.com/ngicks/go-iterator-helper/hiter/stringsiter"
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
				return stringsiter.StringsChunk(stringsSrc, 3)
			},
			Expected: []string{"foo", "bar", "baz"},
			BreakAt:  2,
		}.Test(t)
	})

	t.Run("divide 9 by 4", func(t *testing.T) {
		testcase.One[string]{
			Seq: func() iter.Seq[string] {
				return stringsiter.StringsChunk(stringsSrc, 4)
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
				return stringsiter.StringsChunk(runesSrc, 2)
			},
			Expected: expected,
			BreakAt:  2,
		}.Test(t)

		assert.Equal(t, runesSrc, stringsiter.StringsCollect(len(runesSrc), slices.Values(expected)))
	})

	t.Run("divide by 0", func(t *testing.T) {
		testcase.One[string]{
			Seq: func() iter.Seq[string] {
				return stringsiter.StringsChunk(runesSrc, 0)
			},
			Expected: nil,
		}.Test(t)
	})

	t.Run("divide by -1", func(t *testing.T) {
		testcase.One[string]{
			Seq: func() iter.Seq[string] {
				return stringsiter.StringsChunk(runesSrc, -1)
			},
			Expected: nil,
		}.Test(t)
	})
}

func TestStringsRuneChunk(t *testing.T) {
	t.Run("divide 9 by 3", func(t *testing.T) {
		testcase.One[string]{
			Seq: func() iter.Seq[string] {
				return stringsiter.StringsRuneChunk(stringsSrc, 3)
			},
			Expected: []string{"foo", "bar", "baz"},
			BreakAt:  2,
		}.Test(t)
	})

	t.Run("divide 9 by 4", func(t *testing.T) {
		testcase.One[string]{
			Seq: func() iter.Seq[string] {
				return stringsiter.StringsRuneChunk(stringsSrc, 4)
			},
			Expected: []string{"foob", "arba", "z"},
			BreakAt:  2,
		}.Test(t)
	})

	t.Run("divide emoji by 2", func(t *testing.T) {
		testcase.One[string]{
			Seq: func() iter.Seq[string] {
				return stringsiter.StringsRuneChunk(runesSrc, 2)
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
				return stringsiter.StringsRuneChunk(runesSrc, 0)
			},
			Expected: nil,
		}.Test(t)
	})

	t.Run("divide emoji by -1", func(t *testing.T) {
		testcase.One[string]{
			Seq: func() iter.Seq[string] {
				return stringsiter.StringsRuneChunk(runesSrc, -1)
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
				return stringsiter.StringsSplitFunc(src, -1, nil)
			},
			Expected: []string{"foo", "bar", "baz"},
			BreakAt:  2,
		}.Test(t)

		testcase.One[string]{
			Seq: func() iter.Seq[string] {
				return stringsiter.StringsSplitFunc("", -1, nil)
			},
			Expected: nil,
		}.Test(t)
	})

	t.Run("limit by n", func(t *testing.T) {
		testcase.One[string]{
			Seq: func() iter.Seq[string] {
				return stringsiter.StringsSplitFunc(src, 1, nil)
			},
			Expected: []string{"foo", "bar\nbaz"},
			BreakAt:  1,
		}.Test(t)
	})

	t.Run("StringsCutNewLine", func(t *testing.T) {
		testcase.One[string]{
			Seq: func() iter.Seq[string] {
				return stringsiter.StringsSplitFunc(src, -1, stringsiter.StringsCutNewLine)
			},
			Seqs: []func() iter.Seq[string]{
				func() iter.Seq[string] {
					return stringsiter.StringsSplitFunc(srcCr, -1, stringsiter.StringsCutNewLine)
				},
			},
			Expected: []string{"foo", "bar", "baz"},
			BreakAt:  2,
		}.Test(t)

		testcase.One[string]{
			Seq: func() iter.Seq[string] {
				return stringsiter.StringsSplitFunc(longSingle, -1, stringsiter.StringsCutNewLine)
			},
			Expected: []string{longSingle},
		}.Test(t)
	})

	t.Run("StringsCutUpperCase", func(t *testing.T) {
		testcase.One[string]{
			Seq: func() iter.Seq[string] {
				return stringsiter.StringsSplitFunc("a", -1, stringsiter.StringsCutUpperCase)
			},
			Expected: []string{"a"},
		}.Test(t)

		testcase.One[string]{
			Seq: func() iter.Seq[string] {
				return stringsiter.StringsSplitFunc(src, -1, stringsiter.StringsCutUpperCase)
			},
			Expected: []string{"foo\nbar\nbaz"},
		}.Test(t)

		testcase.One[string]{
			Seq: func() iter.Seq[string] {
				return stringsiter.StringsSplitFunc(srcCase, -1, stringsiter.StringsCutUpperCase)
			},
			Expected: []string{"New", "Http", "Request"},
		}.Test(t)

		testcase.One[string]{
			Seq: func() iter.Seq[string] {
				return stringsiter.StringsSplitFunc(allCapital, -1, stringsiter.StringsCutUpperCase)
			},
			Expected: []string{"S", "T", "O", "P ", "A", "L", "L ", "C", "A", "P", "I", "T", "A", "L"},
		}.Test(t)
	})

	t.Run("StringsCutWord", func(t *testing.T) {
		testcase.One[string]{
			Seq: func() iter.Seq[string] {
				return stringsiter.StringsSplitFunc(allCapital, -1, stringsiter.StringsCutWord)
			},
			Expected: []string{"STOP", "ALL", "CAPITAL"},
		}.Test(t)
		testcase.One[string]{
			Seq: func() iter.Seq[string] {
				return stringsiter.StringsSplitFunc("AAA\t\tBBB    CCC    ", -1, stringsiter.StringsCutWord)
			},
			Expected: []string{"AAA", "BBB", "CCC"},
		}.Test(t)
	})
}
