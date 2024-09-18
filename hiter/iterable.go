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
