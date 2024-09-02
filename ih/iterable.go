package ih

import (
	"iter"
)

// Iterable wraps basic Iter method.
//
// Iter should always return iterators that yield same set of data.
type Iterable[V any] interface {
	Iter() iter.Seq[V]
}

// Iterable2 wraps basic Iter2 method.
//
// Iter2 should always return iterators that yield same set of data.
type Iterable2[K, V any] interface {
	Iter2() iter.Seq2[K, V]
}

// IntoIterable wraps basic IntoIter2 method.
//
// Calling IntoIter may mutate underlying state.
// Therefore calling the method again may also not yield same result.
type IntoIterable[V any] interface {
	IntoIter() iter.Seq[V]
}

// IntoIterable2 wraps basic IntoIter2 method.
//
// Calling IntoIter2 may mutate underlying state.
// Therefore calling the method again may also not yield same result.
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
