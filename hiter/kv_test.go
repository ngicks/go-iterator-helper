package hiter_test

import (
	"iter"
	"slices"
	"testing"

	"github.com/ngicks/go-iterator-helper/hiter"
)

var (
	appendSeq2TestSeq       = hiter.Pairs(hiter.Range(0, 5), hiter.Range(10, 5))
	appendSeq2TestSeqResult = hiter.KeyValues[int, int]{{0, 10}, {1, 9}, {2, 8}, {3, 7}, {4, 6}}
)

func TestValues2(t *testing.T) {
	s := hiter.Collect2(hiter.Values2(appendSeq2TestSeqResult))
	want := appendSeq2TestSeqResult
	if !slices.Equal(s, want) {
		t.Errorf("got %v, want %v", s, want)
	}
}

func TestAppendSeq2(t *testing.T) {
	s := hiter.AppendSeq2(hiter.KeyValues[int, int]{{0, 1}, {3, 2}}, appendSeq2TestSeq)
	want := append(hiter.KeyValues[int, int]{{0, 1}, {3, 2}}, appendSeq2TestSeqResult...)
	if !slices.Equal(s, want) {
		t.Errorf("got %v, want %v", s, want)
	}
}

func TestCollect2(t *testing.T) {
	s := hiter.Collect2(appendSeq2TestSeq)
	want := appendSeq2TestSeqResult
	if !slices.Equal(s, want) {
		t.Errorf("got %v, want %v", s, want)
	}
}

func TestKeyValues(t *testing.T) {
	testCase2[int, int]{
		Seq: func() iter.Seq2[int, int] {
			return hiter.KeyValues[int, int]{
				{2, 1}, {2, 1}, {0, 4}, {-1, 2},
				{2, 1}, {2, 1}, {1, 4}, {-1, 2},
				{2, 1}, {2, 1}, {2, 6}, {-1, 2},
				{2, 1}, {2, 1}, {3, 9}, {-1, 2},
			}.Iter2()
		},
		Expected: []hiter.KeyValue[int, int]{
			{2, 1}, {2, 1}, {0, 4}, {-1, 2},
			{2, 1}, {2, 1}, {1, 4}, {-1, 2},
			{2, 1}, {2, 1}, {2, 6}, {-1, 2},
			{2, 1}, {2, 1}, {3, 9}, {-1, 2},
		},
		BreakAt: 3,
	}.Test(t)
}
