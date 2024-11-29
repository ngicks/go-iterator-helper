package containeriter

import (
	"container/heap"
	"container/list"
	"container/ring"
	"iter"
)

// AppendHeapSeq appends the values from seq to h and returns the extended heap.
// The complexity is O(n) where n = length of total heap content after all element from seq is appended.
func AppendHeapSeq[H heap.Interface, V any](h H, seq iter.Seq[V]) H {
	for v := range seq {
		h.Push(v)
	}
	heap.Init(h)
	return h
}

// AppendListSeq appends the values from seq to l and returns the extended list.
func AppendListSeq[V any](l *list.List, seq iter.Seq[V]) *list.List {
	if l == nil {
		l = list.New()
	}
	for v := range seq {
		l.PushBack(v)
	}
	return l
}

// Collect collects values from seq into a new list and returns it.
func CollectList[V any](seq iter.Seq[V]) *list.List {
	return AppendListSeq(nil, seq)
}

// AppendRingSeq appends the values from seq to r and returns the extended ring.
func AppendRingSeq[V any](r *ring.Ring, seq iter.Seq[V]) *ring.Ring {
	for v := range seq {
		if r == nil {
			r = ring.New(1)
			r.Value = v
			continue
		}
		rr := ring.New(1)
		rr.Value = v
		r = r.Link(rr).Prev()
	}
	return r.Next()
}

// Collect collects values from seq into a new ring and returns it.
func CollectRing[V any](seq iter.Seq[V]) *ring.Ring {
	return AppendRingSeq(nil, seq)
}
