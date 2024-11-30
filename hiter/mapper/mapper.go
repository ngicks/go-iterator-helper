// package sh defines some short hand iterator adapters.
// sh only holds functions which only combine other elements in this module.
package mapper

import (
	"iter"
	"slices"
)

func mapIter[V1, V2 any](fn func(V1) V2, seq iter.Seq[V1]) iter.Seq[V2] {
	return func(yield func(V2) bool) {
		for v1 := range seq {
			if !yield(fn(v1)) {
				return
			}
		}
	}
}

func mapIter2[K1, K2, V1, V2 any](fn func(K1, V1) (K2, V2), seq iter.Seq2[K1, V1]) iter.Seq2[K2, V2] {
	return func(yield func(K2, V2) bool) {
		for k1, v1 := range seq {
			if !yield(fn(k1, v1)) {
				return
			}
		}
	}
}

func filter2[K, V any](f func(K, V) bool, seq iter.Seq2[K, V]) iter.Seq2[K, V] {
	return func(yield func(K, V) bool) {
		for k, v := range seq {
			if f(k, v) && !yield(k, v) {
				return
			}
		}
	}
}

// Collect maps seq by [slices.Collect].
func Collect[V any](seq iter.Seq[iter.Seq[V]]) iter.Seq[[]V] {
	return mapIter(slices.Collect, seq)
}

// Collect2 maps seq by [slices.Collect].
func Collect2[K, V any](seq iter.Seq2[iter.Seq[K], iter.Seq[V]]) iter.Seq2[[]K, []V] {
	return mapIter2(
		func(k iter.Seq[K], v iter.Seq[V]) ([]K, []V) {
			return slices.Collect(k), slices.Collect(v)
		},
		seq,
	)
}

// Clone maps seq by [slices.Clone].
func Clone[S ~[]E, E any](seq iter.Seq[S]) iter.Seq[S] {
	return mapIter(slices.Clone, seq)
}

// Clone2 maps seq by [slices.Clone].
func Clone2[S1 ~[]E1, S2 ~[]E2, E1, E2 any](seq iter.Seq2[S1, S2]) iter.Seq2[S1, S2] {
	return mapIter2(
		func(s1 S1, s2 S2) (S1, S2) { return slices.Clone(s1), slices.Clone(s2) },
		seq,
	)
}
