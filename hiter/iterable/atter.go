package iterable

import (
	"iter"

	"github.com/ngicks/go-iterator-helper/hiter"
)

type Atter[A hiter.Atter[T], T any] struct {
	hiter.Atter[T]
	Indices hiter.Iterable[int]
}

func (a Atter[A, T]) Iter2() iter.Seq2[int, T] {
	return hiter.AtterIndices(a.Atter, a.Indices.Iter())
}
