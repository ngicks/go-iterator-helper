package iterable

import (
	"iter"

	"github.com/ngicks/go-iterator-helper/hiter"
)

type IndexAccessible[A hiter.Atter[T], T any] struct {
	hiter.Atter[T]
	Indices hiter.Iterable[int]
}

func (a IndexAccessible[A, T]) Iter2() iter.Seq2[int, T] {
	return hiter.IndexAccessible(a.Atter, a.Indices.Iter())
}
