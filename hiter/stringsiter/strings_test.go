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
				return stringsiter.Chunk(stringsSrc, 3)
			},
			Expected: []string{"foo", "bar", "baz"},
			BreakAt:  2,
		}.Test(t)
	})

	t.Run("divide 9 by 4", func(t *testing.T) {
		testcase.One[string]{
			Seq: func() iter.Seq[string] {
				return stringsiter.Chunk(stringsSrc, 4)
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
				return stringsiter.Chunk(runesSrc, 2)
			},
			Expected: expected,
			BreakAt:  2,
		}.Test(t)

		assert.Equal(t, runesSrc, stringsiter.Collect(len(runesSrc), slices.Values(expected)))
	})

	t.Run("divide by 0", func(t *testing.T) {
		testcase.One[string]{
			Seq: func() iter.Seq[string] {
				return stringsiter.Chunk(runesSrc, 0)
			},
			Expected: nil,
		}.Test(t)
	})

	t.Run("divide by -1", func(t *testing.T) {
		testcase.One[string]{
			Seq: func() iter.Seq[string] {
				return stringsiter.Chunk(runesSrc, -1)
			},
			Expected: nil,
		}.Test(t)
	})
}

func TestStringsRuneChunk(t *testing.T) {
	t.Run("divide 9 by 3", func(t *testing.T) {
		testcase.One[string]{
			Seq: func() iter.Seq[string] {
				return stringsiter.RuneChunk(stringsSrc, 3)
			},
			Expected: []string{"foo", "bar", "baz"},
			BreakAt:  2,
		}.Test(t)
	})

	t.Run("divide 9 by 4", func(t *testing.T) {
		testcase.One[string]{
			Seq: func() iter.Seq[string] {
				return stringsiter.RuneChunk(stringsSrc, 4)
			},
			Expected: []string{"foob", "arba", "z"},
			BreakAt:  2,
		}.Test(t)
	})

	t.Run("divide emoji by 2", func(t *testing.T) {
		testcase.One[string]{
			Seq: func() iter.Seq[string] {
				return stringsiter.RuneChunk(runesSrc, 2)
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
				return stringsiter.RuneChunk(runesSrc, 0)
			},
			Expected: nil,
		}.Test(t)
	})

	t.Run("divide emoji by -1", func(t *testing.T) {
		testcase.One[string]{
			Seq: func() iter.Seq[string] {
				return stringsiter.RuneChunk(runesSrc, -1)
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
				return stringsiter.SplitFunc(src, -1, nil)
			},
			Expected: []string{"foo", "bar", "baz"},
			BreakAt:  2,
		}.Test(t)

		testcase.One[string]{
			Seq: func() iter.Seq[string] {
				return stringsiter.SplitFunc("", -1, nil)
			},
			Expected: nil,
		}.Test(t)
	})

	t.Run("limit by n", func(t *testing.T) {
		testcase.One[string]{
			Seq: func() iter.Seq[string] {
				return stringsiter.SplitFunc(src, 1, nil)
			},
			Expected: []string{"foo", "bar\nbaz"},
			BreakAt:  1,
		}.Test(t)
	})

	t.Run("StringsCutNewLine", func(t *testing.T) {
		testcase.One[string]{
			Seq: func() iter.Seq[string] {
				return stringsiter.SplitFunc(src, -1, stringsiter.CutNewLine)
			},
			Seqs: []func() iter.Seq[string]{
				func() iter.Seq[string] {
					return stringsiter.SplitFunc(srcCr, -1, stringsiter.CutNewLine)
				},
			},
			Expected: []string{"foo", "bar", "baz"},
			BreakAt:  2,
		}.Test(t)

		testcase.One[string]{
			Seq: func() iter.Seq[string] {
				return stringsiter.SplitFunc(longSingle, -1, stringsiter.CutNewLine)
			},
			Expected: []string{longSingle},
		}.Test(t)
	})

	t.Run("StringsCutUpperCase", func(t *testing.T) {
		testcase.One[string]{
			Seq: func() iter.Seq[string] {
				return stringsiter.SplitFunc("a", -1, stringsiter.CutUpperCase)
			},
			Expected: []string{"a"},
		}.Test(t)

		testcase.One[string]{
			Seq: func() iter.Seq[string] {
				return stringsiter.SplitFunc(src, -1, stringsiter.CutUpperCase)
			},
			Expected: []string{"foo\nbar\nbaz"},
		}.Test(t)

		testcase.One[string]{
			Seq: func() iter.Seq[string] {
				return stringsiter.SplitFunc(srcCase, -1, stringsiter.CutUpperCase)
			},
			Expected: []string{"New", "Http", "Request"},
		}.Test(t)

		testcase.One[string]{
			Seq: func() iter.Seq[string] {
				return stringsiter.SplitFunc(allCapital, -1, stringsiter.CutUpperCase)
			},
			Expected: []string{"S", "T", "O", "P ", "A", "L", "L ", "C", "A", "P", "I", "T", "A", "L"},
		}.Test(t)
	})

	t.Run("StringsCutWord", func(t *testing.T) {
		testcase.One[string]{
			Seq: func() iter.Seq[string] {
				return stringsiter.SplitFunc(allCapital, -1, stringsiter.CutWord)
			},
			Expected: []string{"STOP", "ALL", "CAPITAL"},
		}.Test(t)
		testcase.One[string]{
			Seq: func() iter.Seq[string] {
				return stringsiter.SplitFunc("AAA\t\tBBB    CCC    ", -1, stringsiter.CutWord)
			},
			Expected: []string{"AAA", "BBB", "CCC"},
		}.Test(t)
	})
}
