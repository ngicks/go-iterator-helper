package containeriter_test

import (
	"container/list"
	"container/ring"
	"slices"
	"testing"

	"github.com/ngicks/go-iterator-helper/hiter"
	"github.com/ngicks/go-iterator-helper/hiter/containeriter"
	"gotest.tools/v3/assert"
)

func TestContainerCollect(t *testing.T) {
	seq := hiter.Range(0, 10)

	h1 := &sliceHeap{}
	h2 := containeriter.PushHeapSeq(h1, seq)
	assert.Assert(t, h1 == h2)
	assert.DeepEqual(t, slices.Collect(seq), slices.Collect(containeriter.Heap[int](h2)))

	l1 := list.New()
	for i := range hiter.Range(0, 3) {
		l1.PushBack(i)
	}
	l2 := containeriter.PushBackListSeq(l1, hiter.Range(3, 10))
	assert.Assert(t, l1 == l2)
	assert.DeepEqual(t, slices.Collect(seq), slices.Collect(containeriter.ListAll[int](l2)))
	l3 := containeriter.CollectList(hiter.Range(3, 10))
	assert.Assert(t, l3 != nil)
	assert.DeepEqual(t, slices.Collect(hiter.Range(3, 10)), slices.Collect(containeriter.ListAll[int](l3)))

	r1 := ring.New(3)
	for i := range hiter.Range(0, 3) {
		r1.Value = i
		r1 = r1.Next()
	}
	tail := r1.Prev() // emulate push back, assuming r1 is "root node".
	r2 := containeriter.LinkRingSeq(tail, hiter.Range(3, 10))
	assert.Assert(t, tail == r2)
	assert.DeepEqual(t, slices.Collect(seq), slices.Collect(containeriter.RingAll[int](r2.Prev().Prev())))
	r3 := containeriter.CollectRing(hiter.Range(3, 10))
	assert.Assert(t, r3 != nil)
	assert.DeepEqual(t, slices.Collect(hiter.Range(3, 10)), slices.Collect(containeriter.RingAll[int](r3)))
}

// TestPushFrontListSeq tests the PushFrontListSeq function
func TestPushFrontListSeq(t *testing.T) {
	t.Run("push front to existing list", func(t *testing.T) {
		l := list.New()
		// Pre-populate with some values
		l.PushBack(10)
		l.PushBack(11)

		// Push front values 0, 1, 2 (they will be in reverse order at front)
		result := containeriter.PushFrontListSeq(l, hiter.Range(0, 3))
		assert.Assert(t, result == l, "Should return the same list")

		// Expected order: [2, 1, 0, 10, 11]
		// (Range(0,3) yields 0,1,2, but PushFront means last pushed is first)
		expected := []int{2, 1, 0, 10, 11}
		actual := slices.Collect(containeriter.ListAll[int](result))
		assert.DeepEqual(t, expected, actual)
	})

	t.Run("push front to nil list", func(t *testing.T) {
		result := containeriter.PushFrontListSeq(nil, hiter.Range(5, 3))
		assert.Assert(t, result != nil, "Should create new list when nil passed")

		// Range(5,3) yields 5,4 (descending), PushFront means [4, 5] order
		expected := []int{4, 5}
		actual := slices.Collect(containeriter.ListAll[int](result))
		assert.DeepEqual(t, expected, actual)
	})

	t.Run("push front empty sequence", func(t *testing.T) {
		l := list.New()
		l.PushBack(42)

		result := containeriter.PushFrontListSeq(l, hiter.Range(0, 0)) // empty range
		assert.Assert(t, result == l, "Should return the same list")

		// Should be unchanged
		expected := []int{42}
		actual := slices.Collect(containeriter.ListAll[int](result))
		assert.DeepEqual(t, expected, actual)
	})
}
