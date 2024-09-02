package ih_test

import (
	"iter"
	"slices"
	"testing"

	"github.com/ngicks/go-iterator-helper/ih"
	"gotest.tools/v3/assert"
)

var (
	srcInt1 = []int{12, 76, 8, 9, 923}
	srcInt2 = []int{567, 2, 8, 0, 3}
)

func TestEnumerate(t *testing.T) {
	testCase2[int, int]{
		Seq: func() iter.Seq2[int, int] {
			return ih.Enumerate(slices.Values(srcInt1))
		},
		BreakAt:  2,
		Expected: []ih.KeyValue[int, int]{{0, 12}, {1, 76}, {2, 8}, {3, 9}, {4, 923}},
	}.Test(t)
}

func TestCombine(t *testing.T) {
	testCase2[int, int]{
		Seq: func() iter.Seq2[int, int] {
			return ih.Combine(slices.Values(srcInt1), slices.Values(srcInt2))
		},
		BreakAt:  2,
		Expected: []ih.KeyValue[int, int]{{12, 567}, {76, 2}, {8, 8}, {9, 0}, {923, 3}},
	}.Test(t)
	testCase2[int, int]{
		Seq: func() iter.Seq2[int, int] {
			return ih.Combine(slices.Values(srcInt1[:len(srcInt1)-1]), slices.Values(srcInt2))
		},
		BreakAt:  2,
		Expected: []ih.KeyValue[int, int]{{12, 567}, {76, 2}, {8, 8}, {9, 0}},
	}.Test(t)
	testCase2[int, int]{
		Seq: func() iter.Seq2[int, int] {
			return ih.Combine(slices.Values(srcInt1), slices.Values(srcInt2[:len(srcInt2)-1]))
		},
		BreakAt:  2,
		Expected: []ih.KeyValue[int, int]{{12, 567}, {76, 2}, {8, 8}, {9, 0}},
	}.Test(t)
}

func TestTranspose(t *testing.T) {
	testCase2[int, int]{
		Seq: func() iter.Seq2[int, int] {
			return ih.Transpose(ih.Combine(slices.Values(srcInt1[:len(srcInt1)-1]), slices.Values(srcInt2)))
		},
		BreakAt:  2,
		Expected: []ih.KeyValue[int, int]{{567, 12}, {2, 76}, {8, 8}, {0, 9}},
	}.Test(t)
}

func TestOmitL(t *testing.T) {
	testCase1[int]{
		Seq: func() iter.Seq[int] {
			return ih.OmitL(ih.Combine(slices.Values(srcInt1), slices.Values(srcInt2)))
		},
		BreakAt:  2,
		Expected: []int{12, 76, 8, 9, 923},
	}.Test(t)
}

func TestOmitF(t *testing.T) {
	testCase1[int]{
		Seq: func() iter.Seq[int] {
			return ih.OmitF(ih.Combine(slices.Values(srcInt1), slices.Values(srcInt2)))
		},
		BreakAt:  2,
		Expected: []int{567, 2, 8, 0, 3},
	}.Test(t)
}

func TestOmit(t *testing.T) {
	var i int
	for range ih.Omit(slices.Values(srcInt1)) {
		i++
	}
	assert.Assert(t, len(srcInt1) == i)
}

func TestOmit2(t *testing.T) {
	var i int
	for range ih.Omit2(slices.All(srcInt1)) {
		i++
	}
	assert.Assert(t, len(srcInt1) == i)
}
