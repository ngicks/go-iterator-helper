package iterable

import (
	"iter"

	"github.com/ngicks/go-iterator-helper/hiter"
)

var _ hiter.Iterable[int] = Range[int]{}

// Range adds Iter method to the pair of 2 Numeric values.
type Range[T hiter.Numeric] struct {
	Start, End T
}

// Iter returns an iterator that yields sequential Numeric values in the interval [Start, End).
// Values start from `start` and steps toward `end`.
// At each step value is increased by 1 if start < end, otherwise decreased by 1.
func (r Range[T]) Iter() iter.Seq[T] {
	return hiter.Range(r.Start, r.End)
}
