package iterable_test

import (
	"context"
	"slices"
	"testing"

	"github.com/ngicks/go-iterator-helper/hiter"
	"github.com/ngicks/go-iterator-helper/hiter/iterable"
	"gotest.tools/v3/assert"
)

func TestSliceAll(t *testing.T) {
	data := iterable.SliceAll[string]([]string{"foo", "bar", "baz"})

	// Test Iter method
	result1 := slices.Collect(data.Iter())
	assert.DeepEqual(t, []string{"foo", "bar", "baz"}, result1)

	// Test Iter2 method
	result2 := hiter.Collect2(data.Iter2())
	expected2 := []hiter.KeyValue[int, string]{
		{K: 0, V: "foo"},
		{K: 1, V: "bar"},
		{K: 2, V: "baz"},
	}
	assert.DeepEqual(t, expected2, result2)
}

func TestSliceBackward(t *testing.T) {
	data := iterable.SliceBackward[string]([]string{"foo", "bar", "baz"})

	// Test Iter method
	result1 := slices.Collect(data.Iter())
	assert.DeepEqual(t, []string{"baz", "bar", "foo"}, result1)

	// Test Iter2 method
	result2 := hiter.Collect2(data.Iter2())
	expected2 := []hiter.KeyValue[int, string]{
		{K: 2, V: "baz"},
		{K: 1, V: "bar"},
		{K: 0, V: "foo"},
	}
	assert.DeepEqual(t, expected2, result2)
}

func TestChan(t *testing.T) {
	ch := make(chan int, 3)
	ch <- 1
	ch <- 2
	ch <- 3
	close(ch)

	chanIterable := iterable.Chan[int]{
		Ctx: context.Background(),
		C:   ch,
	}

	result := slices.Collect(chanIterable.IntoIter())
	expected := []int{1, 2, 3}
	assert.DeepEqual(t, expected, result)
}

func TestMapAll(t *testing.T) {
	m := map[string]int{
		"foo": 1,
		"bar": 2,
	}

	mapAll := iterable.MapAll[string, int](m)
	result := hiter.Collect2(mapAll.Iter2())

	// Since map iteration order is not guaranteed, check that we have the right elements
	assert.Assert(t, len(result) == 2)

	found := make(map[string]int)
	for _, kv := range result {
		found[kv.K] = kv.V
	}
	assert.DeepEqual(t, m, found)
}

func TestRange(t *testing.T) {
	rangeIterable := iterable.Range[int]{
		Start: 0,
		End:   5,
	}

	result := slices.Collect(rangeIterable.Iter())
	expected := []int{0, 1, 2, 3, 4}
	assert.DeepEqual(t, expected, result)
}
