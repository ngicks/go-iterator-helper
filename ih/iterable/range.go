package iterable

import (
	"iter"

	"github.com/ngicks/go-iterator-helper/ih"
)

var _ ih.Iterable[int] = Range[int]{}

// Range adds Iter method to the pair of 2 Numeric values.
type Range[T ih.Numeric] struct {
	Start, End T
}

// Iter returns an iterator that yields sequential Numeric values in range [Start, End).
// Values start from `start` and steps toward `end` 1 by 1,
// increased or decreased depending on start < end or not.
func (r Range[T]) Iter() iter.Seq[T] {
	return ih.Range(r.Start, r.End)
}
