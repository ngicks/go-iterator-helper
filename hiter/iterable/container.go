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
	_ hiter.Iterable[any] = List[any]{}
	_ hiter.Iterable[any] = ListBackward[any]{}
)

// List adds Iter method to *list.List.
type List[T any] struct {
	*list.List
}

func (l List[T]) Iter() iter.Seq[T] {
	return hiter.ListAll[T](l.List)
}

// ListBackward adds Iter method to *list.List.
// Iter returns an iterator that traverses list backward.
type ListBackward[T any] struct {
	*list.List
}

func (l ListBackward[T]) Iter() iter.Seq[T] {
	return hiter.ListBackward[T](l.List)
}

var (
	_ hiter.Iterable[any] = Ring[any]{}
	_ hiter.Iterable[any] = RingBackward[any]{}
)

// Ring adds Iter method to *ring.Ring.
type Ring[T any] struct {
	*ring.Ring
}

func (r Ring[T]) Iter() iter.Seq[T] {
	return hiter.RingAll[T](r.Ring)
}

// Ring adds Iter method to *ring.Ring.
// Iter returns an iterator that traverses ring backward.
type RingBackward[T any] struct {
	*ring.Ring
}

func (r RingBackward[T]) Iter() iter.Seq[T] {
	return hiter.RingBackward[T](r.Ring)
}
