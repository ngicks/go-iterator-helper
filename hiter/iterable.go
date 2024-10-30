package hiter

import (
	"iter"
)

// Iterable wraps basic Iter method.
//
// Iter should always return pure / stateless iterators, which always generates same set of data.
type Iterable[V any] interface {
	Iter() iter.Seq[V]
}

// Iterable2 wraps basic Iter2 method.
//
// Iter2 should always return pure / stateless iterators, which always generates same set of pairs.
type Iterable2[K, V any] interface {
	Iter2() iter.Seq2[K, V]
}

// IntoIterable wraps basic IntoIter2 method.
//
// IntoIter might return non-pure / stateful iterators, which would also mutate internal state of implementation.
// Therefore calling the method or invoking the returned iterator multiple times might yield different data without replaying them.
type IntoIterable[V any] interface {
	IntoIter() iter.Seq[V]
}

// IntoIterable2 wraps basic IntoIter2 method.
//
// IntoIter2 might return non-pure / stateful iterators, which would also mutate internal state of implementation.
// Therefore calling the method or invoking the returned iterator multiple times might yield different data without replaying them.
type IntoIterable2[K, V any] interface {
	IntoIter2() iter.Seq2[K, V]
}

var (
	_ Iterable[any]           = FuncIterable[any](nil)
	_ IntoIterable[any]       = FuncIterable[any](nil)
	_ Iterable2[any, any]     = FuncIterable2[any, any](nil)
	_ IntoIterable2[any, any] = FuncIterable2[any, any](nil)
	_ Iterable2[any, any]     = SeqIterable2[any, any](nil)
	_ IntoIterable2[any, any] = SeqIterable2[any, any](nil)
)

type FuncIterable[V any] func() iter.Seq[V]

func (f FuncIterable[V]) Iter() iter.Seq[V] {
	return f()
}

func (f FuncIterable[V]) IntoIter() iter.Seq[V] {
	return f()
}

type FuncIterable2[K, V any] func() iter.Seq2[K, V]

func (f FuncIterable2[K, V]) Iter2() iter.Seq2[K, V] {
	return f()
}

func (f FuncIterable2[K, V]) IntoIter2() iter.Seq2[K, V] {
	return f()
}

type SeqIterable[V any] iter.Seq[V]

// WrapSeqIterable wraps seq into SeqIterable[V].
//
// This is here only for type inference; less typing, easy auto-fill.
func WrapSeqIterable[V any](seq iter.Seq[V]) SeqIterable[V] {
	return SeqIterable[V](seq)
}

func (f SeqIterable[V]) Iter() iter.Seq[V] {
	return iter.Seq[V](f)
}

func (f SeqIterable[V]) IntoIter() iter.Seq[V] {
	return iter.Seq[V](f)
}

type SeqIterable2[K, V any] iter.Seq2[K, V]

// WrapSeqIterable2 wraps seq into SeqIterable2[V].
//
// This is here only for type inference; less typing, easy auto-fill.
func WrapSeqIterable2[K, V any](seq iter.Seq2[K, V]) SeqIterable2[K, V] {
	return SeqIterable2[K, V](seq)
}

func (f SeqIterable2[K, V]) Iter2() iter.Seq2[K, V] {
	return iter.Seq2[K, V](f)
}

func (f SeqIterable2[K, V]) IntoIter2() iter.Seq2[K, V] {
	return iter.Seq2[K, V](f)
}
