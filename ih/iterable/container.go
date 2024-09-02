package iterable

import (
	"container/heap"
	"container/list"
	"container/ring"
	"iter"

	"github.com/ngicks/go-iterator-helper/ih"
)

var _ ih.IntoIterable[any] = Heap[any]{}

// Heap adds IntoIter method to a heap.Interface.
type Heap[T any] struct {
	heap.Interface
}

func (h Heap[T]) IntoIter() iter.Seq[T] {
	return ih.Heap[T](h.Interface)
}

var (
	_ ih.Iterable[any] = List[any]{}
	_ ih.Iterable[any] = ListBackward[any]{}
)

// List adds Iter method to *list.List.
type List[T any] struct {
	*list.List
}

func (l List[T]) Iter() iter.Seq[T] {
	return ih.ListAll[T](l.List)
}

// ListBackward adds Iter method to *list.List.
// Iter returns an iterator that traverses list backward.
type ListBackward[T any] struct {
	*list.List
}

func (l ListBackward[T]) Iter() iter.Seq[T] {
	return ih.ListBackward[T](l.List)
}

var (
	_ ih.Iterable[any] = Ring[any]{}
	_ ih.Iterable[any] = RingBackward[any]{}
)

// Ring adds Iter method to *ring.Ring.
type Ring[T any] struct {
	*ring.Ring
}

func (r Ring[T]) Iter() iter.Seq[T] {
	return ih.RingAll[T](r.Ring)
}

// Ring adds Iter method to *ring.Ring.
// Iter returns an iterator that traverses ring backward.
type RingBackward[T any] struct {
	*ring.Ring
}

func (r RingBackward[T]) Iter() iter.Seq[T] {
	return ih.RingBackward[T](r.Ring)
}
