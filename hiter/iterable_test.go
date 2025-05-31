package hiter_test

import (
	"iter"
	"slices"
	"testing"

	"github.com/ngicks/go-iterator-helper/hiter"
	"gotest.tools/v3/assert"
)

// TestSeqIterable tests the SeqIterable wrapper and its methods
func TestSeqIterable(t *testing.T) {
	seq := slices.Values([]int{1, 2, 3})
	wrapped := hiter.WrapSeqIterable(seq)

	// Test Iter method
	result1 := slices.Collect(wrapped.Iter())
	expected := []int{1, 2, 3}
	assert.DeepEqual(t, expected, result1)

	// Test IntoIter method
	result2 := slices.Collect(wrapped.IntoIter())
	assert.DeepEqual(t, expected, result2)
}

// TestSeqIterable2 tests the SeqIterable2 wrapper and its methods
func TestSeqIterable2(t *testing.T) {
	kvs := []hiter.KeyValue[string, int]{
		{K: "a", V: 1},
		{K: "b", V: 2},
		{K: "c", V: 3},
	}
	seq2 := hiter.Values2(kvs)
	wrapped := hiter.WrapSeqIterable2(seq2)

	// Test Iter2 method
	result1 := hiter.Collect2(wrapped.Iter2())
	assert.DeepEqual(t, kvs, result1)

	// Test IntoIter2 method
	result2 := hiter.Collect2(wrapped.IntoIter2())
	assert.DeepEqual(t, kvs, result2)
}

// TestFuncIterable tests FuncIterable functionality
func TestFuncIterable(t *testing.T) {
	values := []int{10, 20, 30}
	funcIter := hiter.FuncIterable[int](func() iter.Seq[int] {
		return slices.Values(values)
	})

	// Test Iter method
	result1 := slices.Collect(funcIter.Iter())
	assert.DeepEqual(t, values, result1)

	// Test IntoIter method
	result2 := slices.Collect(funcIter.IntoIter())
	assert.DeepEqual(t, values, result2)
}

// TestFuncIterable2 tests FuncIterable2 functionality
func TestFuncIterable2(t *testing.T) {
	kvs := []hiter.KeyValue[string, int]{
		{K: "x", V: 100},
		{K: "y", V: 200},
	}
	funcIter2 := hiter.FuncIterable2[string, int](func() iter.Seq2[string, int] {
		return hiter.Values2(kvs)
	})

	// Test Iter2 method
	result1 := hiter.Collect2(funcIter2.Iter2())
	assert.DeepEqual(t, kvs, result1)

	// Test IntoIter2 method
	result2 := hiter.Collect2(funcIter2.IntoIter2())
	assert.DeepEqual(t, kvs, result2)
}
