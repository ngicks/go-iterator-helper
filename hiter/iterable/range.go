package iterable

import (
	"iter"

	"github.com/ngicks/go-iterator-helper/hiter"
)

var _ hiter.Iterable[int] = Range[int]{}

// Range adds Iter method to the pair of 2 Numeric values.
type Range[T hiter.Numeric] struct {
	Start, End                 T
	ExcludesStart, IncludesEnd bool
}

// Iter returns an iterator that yields sequential Numeric values.
// Values start from `start` and steps toward `end`.
// At each step value is increased by 1 if start < end, otherwise decreased by 1.
// By default, values are in the half open interval [Start, End).
func (r Range[T]) Iter() iter.Seq[T] {
	return hiter.RangeInclusive(r.Start, r.End, !r.ExcludesStart, r.IncludesEnd)
}

func (r Range[T]) Reverse() Range[T] {
	return Range[T]{
		Start:         r.End,
		End:           r.Start,
		ExcludesStart: !r.IncludesEnd,
		IncludesEnd:   !r.ExcludesStart,
	}
}
