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

func (r Range[T]) Reverse() Range[T] {
	// If r was [2, 7), 2 to 6.
	// The reversed r is [6, 1), 6 to 2.
	var start, end T
	if r.Start < r.End {
		start = r.End - 1
		if start > r.End {
			panic("underflow")
		}
		end = r.Start - 1
		if end > r.Start {
			panic("underflow")
		}
	} else {
		start = r.End + 1
		if start < r.End {
			panic("overflow")
		}
		end = r.Start + 1
		if end < r.Start {
			panic("overflow")
		}
	}
	return Range[T]{
		Start: start,
		End:   end,
	}
}
