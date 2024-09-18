package iterable

import (
	"container/heap"
	"container/list"
	"container/ring"
	"iter"

	"github.com/ngicks/go-iterator-helper/hiter"
)

var _ hiter.IntoIterable[any] = Heap[any]{}

// Heap adds IntoIter method to a heap.Interface.
type Heap[T any] struct {
	heap.Interface
}

func (h Heap[T]) IntoIter() iter.Seq[T] {
	return hiter.Heap[T](h.Interface)
}

var (
	_ hiter.Iterable[any] = ListAll[any]{}
	_ hiter.Iterable[any] = ListBackward[any]{}
)

// ListAll adds Iter method to *list.ListAll.
type ListAll[T any] struct {
	*list.List
}

func (l ListAll[T]) Iter() iter.Seq[T] {
	return hiter.ListAll[T](l.List)
}

// ListElementAll adds Iter method to *list.Element.
type ListElementAll[T any] struct {
	*list.Element
}

func (l ListElementAll[T]) Iter() iter.Seq[T] {
	return hiter.ListElementAll[T](l.Element)
}

// ListBackward adds Iter method to *list.List.
// Iter returns an iterator that traverses list backward.
type ListBackward[T any] struct {
	*list.List
}

func (l ListBackward[T]) Iter() iter.Seq[T] {
	return hiter.ListBackward[T](l.List)
}

// ListElementBackward adds Iter method to *list.Element.
// Iter returns an iterator that traverses list backward.
type ListElementBackward[T any] struct {
	*list.Element
}

func (l ListElementBackward[T]) Iter() iter.Seq[T] {
	return hiter.ListElementBackward[T](l.Element)
}

var (
	_ hiter.Iterable[any] = RingAll[any]{}
	_ hiter.Iterable[any] = RingBackward[any]{}
)

// RingAll adds Iter method to *ring.RingAll.
type RingAll[T any] struct {
	*ring.Ring
}

func (r RingAll[T]) Iter() iter.Seq[T] {
	return hiter.RingAll[T](r.Ring)
}

// RingBackward adds Iter method to *ring.Ring.
// Iter returns an iterator that traverses ring backward.
type RingBackward[T any] struct {
	*ring.Ring
}

func (r RingBackward[T]) Iter() iter.Seq[T] {
	return hiter.RingBackward[T](r.Ring)
}
