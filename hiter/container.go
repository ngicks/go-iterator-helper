package hiter

import (
	"container/heap"
	"container/list"
	"container/ring"
	"iter"
)

// Heap returns an iterator over heap.Interface.
// Consuming iter.Seq[T] also consumes h.
// To avoid this, the caller must clone input h before passing to Heap.
func Heap[T any](h heap.Interface) iter.Seq[T] {
	return func(yield func(T) bool) {
		for h.Len() > 0 {
			popped := heap.Pop(h)
			if !yield(popped.(T)) {
				return
			}
		}
	}
}

// ListAll returns an iterator over l.
func ListAll[T any](l *list.List) iter.Seq[T] {
	return ListElementAll[T](l.Front())
}

// ListElementAll returns an iterator over from ele to end of the list.
func ListElementAll[T any](ele *list.Element) iter.Seq[T] {
	return func(yield func(T) bool) {
		for ; ele != nil; ele = ele.Next() {
			if !yield(ele.Value.(T)) {
				return
			}
		}
	}
}

// ListBackward returns an iterator over l,
// traversing it backward by calling Back and Prev.
func ListBackward[T any](l *list.List) iter.Seq[T] {
	return ListElementBackward[T](l.Back())
}

// ListElementBackward returns an iterator over from ele to start of the list.
func ListElementBackward[T any](ele *list.Element) iter.Seq[T] {
	return func(yield func(T) bool) {
		for ; ele != nil; ele = ele.Prev() {
			if !yield(ele.Value.(T)) {
				return
			}
		}
	}
}

// Ring returns an iterator over r.
// by traversing from r and consecutively calling Next.
func RingAll[T any](r *ring.Ring) iter.Seq[T] {
	return func(yield func(T) bool) {
		if !yield(r.Value.(T)) {
			return
		}
		for n := r.Next(); n != r; n = n.Next() {
			if !yield(n.Value.(T)) {
				return
			}
		}
	}
}

// RingBackward returns an iterator over r,
// traversing it backward starting from r and consecutively calling Prev.
func RingBackward[T any](r *ring.Ring) iter.Seq[T] {
	return func(yield func(T) bool) {
		if !yield(r.Value.(T)) {
			return
		}
		for n := r.Prev(); n != r; n = n.Prev() {
			if !yield(n.Value.(T)) {
				return
			}
		}
	}
}
