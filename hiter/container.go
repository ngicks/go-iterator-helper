package hiter

import (
	"container/heap"
	"container/list"
	"container/ring"
	"iter"
)

// Heap returns an iterator over heap.Interface.
// Consuming iter.Seq[V] also consumes h.
// To avoid this, callers must clone input h before passing to Heap.
func Heap[V any](h heap.Interface) iter.Seq[V] {
	return func(yield func(V) bool) {
		for h.Len() > 0 {
			popped := heap.Pop(h)
			if !yield(popped.(V)) {
				return
			}
		}
	}
}

// ListAll returns an iterator over all element of l starting from l.Front().
// ListAll assumes Values of all element are type V.
// If other than that or nil, the returned iterator may panic on invocation.
func ListAll[V any](l *list.List) iter.Seq[V] {
	return ListElementAll[V](l.Front())
}

// ListElementAll returns an iterator over from ele to end of the list.
// ListElementAll assumes Values of all element are type V.
// If other than that or nil, the returned iterator may panic on invocation.
func ListElementAll[V any](ele *list.Element) iter.Seq[V] {
	return func(yield func(V) bool) {
		// shadowing ele, no state in the seq closure as much as possible.
		for ele := ele; ele != nil; ele = ele.Next() {
			if !yield(ele.Value.(V)) {
				return
			}
		}
	}
}

// ListBackward returns an iterator over all element of l starting from l.Back().
// ListBackward assumes Values of all element are type V.
// If other than that or nil, the returned iterator may panic on invocation.
func ListBackward[V any](l *list.List) iter.Seq[V] {
	return ListElementBackward[V](l.Back())
}

// ListElementBackward returns an iterator over from ele to start of the list.
// ListElementBackward assumes Values of all element are type V.
// If other than that or nil, the returned iterator may panic on invocation.
func ListElementBackward[V any](ele *list.Element) iter.Seq[V] {
	return func(yield func(V) bool) {
		// no state in in the seq closure as much as possible.
		for ele := ele; ele != nil; ele = ele.Prev() {
			if !yield(ele.Value.(V)) {
				return
			}
		}
	}
}

// Ring returns an iterator over r.
// The returned iterator generates data assuming Values of all ring elements are type V.
// It yields r.Value traversing by consecutively calling Next, and stops when it finds r again.
// Removing r from the ring after it started iteration may make it iterate forever.
func RingAll[V any](r *ring.Ring) iter.Seq[V] {
	return func(yield func(V) bool) {
		if !yield(r.Value.(V)) {
			return
		}
		for n := r.Next(); n != r; n = n.Next() {
			if !yield(n.Value.(V)) {
				return
			}
		}
	}
}

// RingBackward returns an iterator over r.
// The returned iterator generates data assuming Values of all ring elements are type V.
// It yields r.Value traversing by consecutively calling Prev, and stops when it finds r again.
// Removing r from the ring after it started iteration may make it iterate forever.
func RingBackward[V any](r *ring.Ring) iter.Seq[V] {
	return func(yield func(V) bool) {
		if !yield(r.Value.(V)) {
			return
		}
		for n := r.Prev(); n != r; n = n.Prev() {
			if !yield(n.Value.(V)) {
				return
			}
		}
	}
}
