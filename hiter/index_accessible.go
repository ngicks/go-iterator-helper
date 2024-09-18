package hiter

import "iter"

type Atter[T any] interface {
	At(i int) T
}

// IndexAccessible returns an iterator over indices and values of a associated to the indices.
// If indices generates an out-of-range index, the behavior is not defined and may differs among Atter implementations.
func IndexAccessible[A Atter[T], T any](a A, indices iter.Seq[int]) iter.Seq2[int, T] {
	return func(yield func(int, T) bool) {
		for i := range indices {
			if !yield(i, a.At(i)) {
				return
			}
		}
	}
}
