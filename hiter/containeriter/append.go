package containeriter

import (
	"container/heap"
	"container/list"
	"container/ring"
	"iter"
)

// PushHeapSeq pushes the values from seq to h and returns the extended heap.
func PushHeapSeq[H heap.Interface, V any](h H, seq iter.Seq[V]) H {
	for v := range seq {
		h.Push(v)
	}
	return h
}

// PushBackListSeq pushes-back the values from seq to l and returns the extended list.
func PushBackListSeq[V any](l *list.List, seq iter.Seq[V]) *list.List {
	if l == nil {
		l = list.New()
	}
	for v := range seq {
		l.PushBack(v)
	}
	return l
}

// PushFrontListSeq pushes-front the values from seq to l and returns the extended list.
func PushFrontListSeq[V any](l *list.List, seq iter.Seq[V]) *list.List {
	if l == nil {
		l = list.New()
	}
	for v := range seq {
		l.PushFront(v)
	}
	return l
}

// Collect collects values from seq into a new list and returns it.
func CollectList[V any](seq iter.Seq[V]) *list.List {
	return PushBackListSeq(nil, seq)
}

// LinkRingSeq links all values from seq to r and returns r.
func LinkRingSeq[V any](r *ring.Ring, seq iter.Seq[V]) *ring.Ring {
	root := r
	for v := range seq {
		if r == nil {
			r = ring.New(1)
			r.Value = v
			root = r
			continue
		}
		next := ring.New(1)
		next.Value = v
		r = r.Link(next).Prev()
	}
	return root
}

// Collect collects values from seq into a new ring and returns it.
func CollectRing[V any](seq iter.Seq[V]) *ring.Ring {
	return LinkRingSeq(nil, seq)
}
