package hiter

import (
	"iter"
)

type Atter[T any] interface {
	At(i int) T
}

// AtterIndices returns an iterator over pairs of indices and values retrieved from an Atter using the indices.
// The indices sequence determines the order of access and must provide valid indices within the bounds of the Atter.
// If an out-of-range index is supplied, the behavior is undefined and may vary based on the Atter implementation.
func AtterIndices[A Atter[T], T any](a A, indices iter.Seq[int]) iter.Seq2[int, T] {
	return func(yield func(int, T) bool) {
		for i := range indices {
			if !yield(i, a.At(i)) {
				return
			}
		}
	}
}

// AtterRange is like [AtterIndices] but indices is result of calling [Range] with start and end.
func AtterRange[A Atter[T], T any](a A, start, end int) iter.Seq2[int, T] {
	return AtterIndices(a, Range(start, end))
}

type Lenner interface {
	Len() int
}

// AtterAll is like [AtterRange] but start is 0 and end is result of [Lenner.Len].
func AtterAll[A interface {
	Atter[T]
	Lenner
}, T any](a A) iter.Seq2[int, T] {
	return AtterRange(a, 0, a.Len())
}
