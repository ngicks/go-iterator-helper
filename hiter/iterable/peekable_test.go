package iterable_test

import (
	"iter"
	"slices"
	"testing"

	"github.com/ngicks/go-iterator-helper/hiter"
	"github.com/ngicks/go-iterator-helper/hiter/iterable"
	"github.com/ngicks/go-iterator-helper/x/exp/xiter"
	"gotest.tools/v3/assert"
	"gotest.tools/v3/assert/cmp"
)

func TestPeekable_resumable(t *testing.T) {
	testResumable(t, func(src iter.Seq[int]) *iterable.Peekable[int] { return iterable.NewPeekable(src) })
}

func TestPeekable2_resumable(t *testing.T) {
	testResumable2(t, func(src iter.Seq2[int, int]) *iterable.Peekable2[int, int] { return iterable.NewPeekable2(src) })
}

func TestPeekable(t *testing.T) {
	res := iterable.NewPeekable(hiter.Range(0, 20))

	assert.DeepEqual(t, slices.Collect(hiter.Range(0, 2)), res.Peek(2))
	assert.DeepEqual(t, slices.Collect(hiter.Range(0, 2)), res.Peek(2))
	assert.DeepEqual(t, slices.Collect(hiter.Range(0, 2)), res.Peek(2))

	assert.DeepEqual(t, slices.Collect(hiter.Range(0, 5)), res.Peek(5))
	assert.DeepEqual(t, slices.Collect(hiter.Range(0, 5)), res.Peek(5))

	assert.DeepEqual(t, slices.Collect(hiter.Range(0, 10)), slices.Collect(xiter.Limit(res.IntoIter(), 10)))
	assert.DeepEqual(t, slices.Collect(hiter.Range(10, 12)), res.Peek(2))

	rest := slices.Collect(hiter.Range(10, 15))
	for i, v := range xiter.Limit2(hiter.Enumerate(res.IntoIter()), 5) {
		res.Peek(3)
		assert.Equal(t, rest[i], v)
	}
	assert.DeepEqual(t, slices.Collect(hiter.Range(15, 17)), res.Peek(2))
	// 3 elements peeked and buffered.
	res.Stop()
	assert.DeepEqual(
		t,
		slices.Collect(hiter.Range(15, 18)),
		slices.Collect(res.IntoIter()),
	)

	assert.Assert(t, cmp.Len(res.Peek(5), 0))
}

func TestPeekable2(t *testing.T) {
	res := iterable.NewPeekable2(hiter.Enumerate(hiter.Range(0, 20)))

	collectRange := func(start, end int) []hiter.KeyValue[int, int] {
		return hiter.Collect2(hiter.Pairs(hiter.Range(start, end), hiter.Range(start, end)))
	}
	assert.DeepEqual(t, collectRange(0, 2), res.Peek(2))
	assert.DeepEqual(t, collectRange(0, 2), res.Peek(2))
	assert.DeepEqual(t, collectRange(0, 2), res.Peek(2))

	assert.DeepEqual(t, collectRange(0, 5), res.Peek(5))
	assert.DeepEqual(t, collectRange(0, 5), res.Peek(5))

	assert.DeepEqual(t, collectRange(0, 10), hiter.Collect2(xiter.Limit2(res.IntoIter2(), 10)))
	assert.DeepEqual(t, collectRange(10, 12), res.Peek(2))

	rest := collectRange(10, 15)
	for i, v := range xiter.Limit2(hiter.Enumerate(hiter.ToKeyValue(res.IntoIter2())), 5) {
		res.Peek(3)
		assert.Equal(t, rest[i], v)
	}
	assert.DeepEqual(t, collectRange(15, 17), res.Peek(2))
	// 3 elements peeked and buffered.
	res.Stop()
	assert.DeepEqual(
		t,
		collectRange(15, 18),
		hiter.Collect2(res.IntoIter2()),
	)

	assert.Assert(t, cmp.Len(res.Peek(5), 0))
}
