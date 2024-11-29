package iterable

import (
	"container/heap"
	"container/list"
	"container/ring"
	"iter"

	"github.com/ngicks/go-iterator-helper/hiter"
	"github.com/ngicks/go-iterator-helper/hiter/containeriter"
)

var _ hiter.IntoIterable[any] = Heap[any]{}

// Heap adds IntoIter method to a heap.Interface.
type Heap[T any] struct {
	heap.Interface
}

func (h Heap[T]) IntoIter() iter.Seq[T] {
	return containeriter.Heap[T](h.Interface)
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
	return containeriter.ListAll[T](l.List)
}

func (l ListAll[T]) Reverse() ListBackward[T] {
	return ListBackward[T](l)
}

// ListElementAll adds Iter method to *list.Element.
type ListElementAll[T any] struct {
	*list.Element
}

func (l ListElementAll[T]) Iter() iter.Seq[T] {
	return containeriter.ListElementAll[T](l.Element)
}

func (l ListElementAll[T]) Reverse() ListElementBackward[T] {
	return ListElementBackward[T](l)
}

// ListBackward adds Iter method to *list.List.
// Iter returns an iterator that traverses list backward.
type ListBackward[T any] struct {
	*list.List
}

func (l ListBackward[T]) Iter() iter.Seq[T] {
	return containeriter.ListBackward[T](l.List)
}

func (l ListBackward[T]) Reverse() ListAll[T] {
	return ListAll[T](l)
}

// ListElementBackward adds Iter method to *list.Element.
// Iter returns an iterator that traverses list backward.
type ListElementBackward[T any] struct {
	*list.Element
}

func (l ListElementBackward[T]) Iter() iter.Seq[T] {
	return containeriter.ListElementBackward[T](l.Element)
}

func (l ListElementBackward[T]) Reverse() ListElementAll[T] {
	return ListElementAll[T](l)
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
	return containeriter.RingAll[T](r.Ring)
}

func (r RingAll[T]) Reverse() RingBackward[T] {
	return RingBackward[T]{r.Ring.Prev()}
}

// RingBackward adds Iter method to *ring.Ring.
// Iter returns an iterator that traverses ring backward.
type RingBackward[T any] struct {
	*ring.Ring
}

func (r RingBackward[T]) Iter() iter.Seq[T] {
	return containeriter.RingBackward[T](r.Ring)
}

func (r RingBackward[T]) Reverse() RingAll[T] {
	return RingAll[T]{r.Ring.Next()}
}
