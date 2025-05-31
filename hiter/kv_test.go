package hiter_test

import (
	"iter"
	"slices"
	"testing"

	"github.com/ngicks/go-iterator-helper/hiter"
	"github.com/ngicks/go-iterator-helper/hiter/internal/testcase"
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
	s := hiter.AppendSeq2(
		hiter.KeyValues[int, int]{
			{0, 1},
			{3, 2},
		},
		appendSeq2TestSeq,
	)
	want := append(hiter.KeyValues[int, int]{
		{0, 1},
		{3, 2},
	},
		appendSeq2TestSeqResult...,
	)
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
	testcase.Two[int, int]{
		Seq: func() iter.Seq2[int, int] {
			return hiter.KeyValues[int, int]{
				{2, 1},
				{2, 1},
				{0, 4},
				{-1, 2},
				{2, 1},
				{2, 1},
				{1, 4},
				{-1, 2},
				{2, 1},
				{2, 1},
				{2, 6},
				{-1, 2},
				{2, 1},
				{2, 1},
				{3, 9},
				{-1, 2},
			}.Iter2()
		},
		Expected: []hiter.KeyValue[int, int]{
			{2, 1},
			{2, 1},
			{0, 4},
			{-1, 2},
			{2, 1},
			{2, 1},
			{1, 4},
			{-1, 2},
			{2, 1},
			{2, 1},
			{2, 6},
			{-1, 2},
			{2, 1},
			{2, 1},
			{3, 9},
			{-1, 2},
		},
		BreakAt: 3,
	}.Test(t)
}

// TestKVPair tests the KVPair constructor function
func TestKVPair(t *testing.T) {
	kv := hiter.KVPair("key", 42)
	if kv.K != "key" || kv.V != 42 {
		t.Errorf("KVPair failed: got %+v, want {K: \"key\", V: 42}", kv)
	}
}

// TestKeyValueUnpack tests the Unpack method
func TestKeyValueUnpack(t *testing.T) {
	kv := hiter.KeyValue[string, int]{K: "test", V: 123}
	k, v := kv.Unpack()
	if k != "test" || v != 123 {
		t.Errorf("Unpack failed: got (%v, %v), want (\"test\", 123)", k, v)
	}
}

// TestToKeyValue tests converting iter.Seq2 to iter.Seq[KeyValue]
func TestToKeyValue(t *testing.T) {
	seq2 := func(yield func(string, int) bool) {
		if !yield("a", 1) || !yield("b", 2) || !yield("c", 3) {
			return
		}
	}

	kvSeq := hiter.ToKeyValue(seq2)
	result := slices.Collect(kvSeq)

	expected := []hiter.KeyValue[string, int]{
		{K: "a", V: 1},
		{K: "b", V: 2},
		{K: "c", V: 3},
	}

	if !slices.Equal(result, expected) {
		t.Errorf("ToKeyValue failed: got %v, want %v", result, expected)
	}
}

// TestFromKeyValue tests converting iter.Seq[KeyValue] back to iter.Seq2
func TestFromKeyValue(t *testing.T) {
	kvs := []hiter.KeyValue[string, int]{
		{K: "x", V: 10},
		{K: "y", V: 20},
		{K: "z", V: 30},
	}

	kvSeq := slices.Values(kvs)
	seq2 := hiter.FromKeyValue(kvSeq)
	result := hiter.Collect2(seq2)

	if !slices.Equal(result, kvs) {
		t.Errorf("FromKeyValue failed: got %v, want %v", result, kvs)
	}
}

// TestToKeyValueFromKeyValueRoundtrip tests that ToKeyValue and FromKeyValue are inverses
func TestToKeyValueFromKeyValueRoundtrip(t *testing.T) {
	seq2 := func(yield func(int, string) bool) {
		if !yield(1, "one") || !yield(2, "two") || !yield(3, "three") {
			return
		}
	}

	// Convert to KeyValue and back
	kvSeq := hiter.ToKeyValue(seq2)
	backToSeq2 := hiter.FromKeyValue(kvSeq)
	result := hiter.Collect2(backToSeq2)

	expected := []hiter.KeyValue[int, string]{
		{K: 1, V: "one"},
		{K: 2, V: "two"},
		{K: 3, V: "three"},
	}

	if !slices.Equal(result, expected) {
		t.Errorf("Roundtrip failed: got %v, want %v", result, expected)
	}
}

// TestToKeyValueEarlyBreak tests that ToKeyValue respects early termination
func TestToKeyValueEarlyBreak(t *testing.T) {
	seq2 := func(yield func(int, int) bool) {
		for i := 0; i < 10; i++ {
			if !yield(i, i*i) {
				return
			}
		}
	}

	kvSeq := hiter.ToKeyValue(seq2)
	var result []hiter.KeyValue[int, int]
	for kv := range kvSeq {
		result = append(result, kv)
		if len(result) == 3 {
			break
		}
	}

	expected := []hiter.KeyValue[int, int]{
		{K: 0, V: 0},
		{K: 1, V: 1},
		{K: 2, V: 4},
	}

	if !slices.Equal(result, expected) {
		t.Errorf("ToKeyValue early break failed: got %v, want %v", result, expected)
	}
}

// TestFromKeyValueEarlyBreak tests that FromKeyValue respects early termination
func TestFromKeyValueEarlyBreak(t *testing.T) {
	kvs := []hiter.KeyValue[string, int]{
		{K: "a", V: 1},
		{K: "b", V: 2},
		{K: "c", V: 3},
		{K: "d", V: 4},
		{K: "e", V: 5},
	}

	kvSeq := slices.Values(kvs)
	seq2 := hiter.FromKeyValue(kvSeq)

	var result []hiter.KeyValue[string, int]
	for k, v := range seq2 {
		result = append(result, hiter.KeyValue[string, int]{K: k, V: v})
		if len(result) == 3 {
			break
		}
	}

	expected := []hiter.KeyValue[string, int]{
		{K: "a", V: 1}, {K: "b", V: 2}, {K: "c", V: 3},
	}

	if !slices.Equal(result, expected) {
		t.Errorf("FromKeyValue early break failed: got %v, want %v", result, expected)
	}
}
