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

func TestStringChunk(t *testing.T) {
	t.Run("divide 9 by 3", func(t *testing.T) {
		testCase1[string]{
			Seq: func() iter.Seq[string] {
				return ih.StringChunk(stringsSrc, 3)
			},
			Expected: []string{"foo", "bar", "baz"},
			BreakAt:  2,
		}.Test(t)
	})

	t.Run("divide 9 by 4", func(t *testing.T) {
		testCase1[string]{
			Seq: func() iter.Seq[string] {
				return ih.StringChunk(stringsSrc, 4)
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
				return ih.StringChunk(runesSrc, 2)
			},
			Expected: expected,
			BreakAt:  2,
		}.Test(t)

		assert.Equal(t, runesSrc, ih.CollectString(slices.Values(expected), len(runesSrc)))
	})

	t.Run("divide by 0", func(t *testing.T) {
		testCase1[string]{
			Seq: func() iter.Seq[string] {
				return ih.StringChunk(runesSrc, 0)
			},
			Expected: nil,
		}.Test(t)
	})

	t.Run("divide by -1", func(t *testing.T) {
		testCase1[string]{
			Seq: func() iter.Seq[string] {
				return ih.StringChunk(runesSrc, -1)
			},
			Expected: nil,
		}.Test(t)
	})
}

func TestStringChunkRune(t *testing.T) {
	t.Run("divide 9 by 3", func(t *testing.T) {
		testCase1[string]{
			Seq: func() iter.Seq[string] {
				return ih.StringRuneChunk(stringsSrc, 3)
			},
			Expected: []string{"foo", "bar", "baz"},
			BreakAt:  2,
		}.Test(t)
	})

	t.Run("divide 9 by 4", func(t *testing.T) {
		testCase1[string]{
			Seq: func() iter.Seq[string] {
				return ih.StringRuneChunk(stringsSrc, 4)
			},
			Expected: []string{"foob", "arba", "z"},
			BreakAt:  2,
		}.Test(t)
	})

	t.Run("divide emoji by 2", func(t *testing.T) {
		testCase1[string]{
			Seq: func() iter.Seq[string] {
				return ih.StringRuneChunk(runesSrc, 2)
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
				return ih.StringRuneChunk(runesSrc, 0)
			},
			Expected: nil,
		}.Test(t)
	})

	t.Run("divide emoji by -1", func(t *testing.T) {
		testCase1[string]{
			Seq: func() iter.Seq[string] {
				return ih.StringRuneChunk(runesSrc, -1)
			},
			Expected: nil,
		}.Test(t)
	})
}
