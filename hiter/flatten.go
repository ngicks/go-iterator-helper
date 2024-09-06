package hiter

import "iter"

// Flatten returns an iterator over slices yielded from seq.
func Flatten[S ~[]E, E any](seq iter.Seq[S]) iter.Seq[E] {
	return func(yield func(E) bool) {
		for s := range seq {
			for _, t := range s {
				if !yield(t) {
					return
				}
			}
		}
	}
}

// FlattenF returns an iterator over pairs of slice and non-slice.
// While iterating over slices, the latter part of pair is repeated.
func FlattenF[S1 ~[]E1, E1 any, E2 any](seq iter.Seq2[S1, E2]) iter.Seq2[E1, E2] {
	return func(yield func(E1, E2) bool) {
		for s, g := range seq {
			for _, e := range s {
				if !yield(e, g) {
					return
				}
			}
		}
	}
}

// FlattenL returns an iterator over pairs of non-slice and slice.
// While iterating over slices, the former part of pair is repeated.
func FlattenL[S2 ~[]E2, E1 any, E2 any](seq iter.Seq2[E1, S2]) iter.Seq2[E1, E2] {
	return func(yield func(E1, E2) bool) {
		for e, t := range seq {
			for _, g := range t {
				if !yield(e, g) {
					return
				}
			}
		}
	}
}
