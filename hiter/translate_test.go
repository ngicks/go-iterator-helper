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
	srcInt1 = []int{12, 76, 8, 9, 923}
	srcInt2 = []int{567, 2, 8, 0, 3}
)

func TestEnumerate(t *testing.T) {
	testcase.TestCase2[int, int]{
		Seq: func() iter.Seq2[int, int] {
			return hiter.Enumerate(slices.Values(srcInt1))
		},
		BreakAt:  2,
		Expected: []hiter.KeyValue[int, int]{{0, 12}, {1, 76}, {2, 8}, {3, 9}, {4, 923}},
	}.Test(t)
}

func TestPairs(t *testing.T) {
	testcase.TestCase2[int, int]{
		Seq: func() iter.Seq2[int, int] {
			return hiter.Pairs(slices.Values(srcInt1), slices.Values(srcInt2))
		},
		BreakAt:  2,
		Expected: []hiter.KeyValue[int, int]{{12, 567}, {76, 2}, {8, 8}, {9, 0}, {923, 3}},
	}.Test(t)
	testcase.TestCase2[int, int]{
		Seq: func() iter.Seq2[int, int] {
			return hiter.Pairs(slices.Values(srcInt1[:len(srcInt1)-1]), slices.Values(srcInt2))
		},
		BreakAt:  2,
		Expected: []hiter.KeyValue[int, int]{{12, 567}, {76, 2}, {8, 8}, {9, 0}},
	}.Test(t)
	testcase.TestCase2[int, int]{
		Seq: func() iter.Seq2[int, int] {
			return hiter.Pairs(slices.Values(srcInt1), slices.Values(srcInt2[:len(srcInt2)-1]))
		},
		BreakAt:  2,
		Expected: []hiter.KeyValue[int, int]{{12, 567}, {76, 2}, {8, 8}, {9, 0}},
	}.Test(t)
}

func TestTranspose(t *testing.T) {
	testcase.TestCase2[int, int]{
		Seq: func() iter.Seq2[int, int] {
			return hiter.Transpose(hiter.Pairs(slices.Values(srcInt1[:len(srcInt1)-1]), slices.Values(srcInt2)))
		},
		BreakAt:  2,
		Expected: []hiter.KeyValue[int, int]{{567, 12}, {2, 76}, {8, 8}, {0, 9}},
	}.Test(t)
}

func TestOmitL(t *testing.T) {
	testcase.TestCase1[int]{
		Seq: func() iter.Seq[int] {
			return hiter.OmitL(hiter.Pairs(slices.Values(srcInt1), slices.Values(srcInt2)))
		},
		BreakAt:  2,
		Expected: []int{12, 76, 8, 9, 923},
	}.Test(t)
}

func TestOmitF(t *testing.T) {
	testcase.TestCase1[int]{
		Seq: func() iter.Seq[int] {
			return hiter.OmitF(hiter.Pairs(slices.Values(srcInt1), slices.Values(srcInt2)))
		},
		BreakAt:  2,
		Expected: []int{567, 2, 8, 0, 3},
	}.Test(t)
}

func TestOmit(t *testing.T) {
	var i int
	for range hiter.Omit(slices.Values(srcInt1)) {
		i++
	}
	assert.Assert(t, len(srcInt1) == i)
}

func TestOmit2(t *testing.T) {
	var i int
	for range hiter.Omit2(slices.All(srcInt1)) {
		i++
	}
	assert.Assert(t, len(srcInt1) == i)
}

func TestUniteBy(t *testing.T) {
	src := hiter.KeyValues[int, string]{{1, "foo"}, {2, "bar"}, {3, "baz"}}
	united := hiter.Unify(
		func(k int, v string) hiter.KeyValue[int, string] { return hiter.KeyValue[int, string]{k, v} },
		src.Iter2(),
	)
	assert.DeepEqual(t, src, slices.AppendSeq[hiter.KeyValues[int, string]](nil, united))
	var mid hiter.KeyValues[int, string]
	for i, kv := range hiter.Enumerate(united) {
		mid = append(mid, kv)
		if i == 1 {
			break
		}
	}
	assert.DeepEqual(t, src[:2], mid)
}

func TestDivideBy(t *testing.T) {
	src := []hiter.KeyValue[int, string]{{1, "foo"}, {2, "bar"}, {3, "baz"}}
	divided := hiter.Divide(
		func(kv hiter.KeyValue[int, string]) (int, string) { return kv.K, kv.V },
		slices.Values(src),
	)
	assert.DeepEqual(t, src, hiter.Collect2(divided))
	var mid []hiter.KeyValue[int, string]
	var i int
	for k, v := range divided {
		mid = append(mid, hiter.KeyValue[int, string]{k, v})
		if i == 1 {
			break
		}
		i++
	}
	assert.DeepEqual(t, src[:2], mid)
}
