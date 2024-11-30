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
