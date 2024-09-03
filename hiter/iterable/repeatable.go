package iterable

import (
	"iter"

	"github.com/ngicks/go-iterator-helper/hiter"
)

var (
	_ hiter.Iterable[int]       = Repeatable[int]{}
	_ hiter.Iterable2[int, int] = Repeatable2[int, int]{}
)

// Repeatable generates an iterator that generates V N times.
type Repeatable[V any] struct {
	V V
	N int
}

func (r Repeatable[V]) Iter() iter.Seq[V] {
	return hiter.Repeat(r.V, r.N)
}

// Repeatable2 generates an iterator that generates pairs of K and V N times.
type Repeatable2[K, V any] struct {
	K K
	V V
	N int
}

func (r Repeatable2[K, V]) Iter2() iter.Seq2[K, V] {
	return hiter.Repeat2(r.K, r.V, r.N)
}
