// package sh defines some short hand iterator adapters.
// sh only holds functions which only combine other elements in this module.
package mapper

import (
	"iter"
	"slices"

	"github.com/ngicks/go-iterator-helper/hiter/internal/adapter"
)

// Collect maps seq by [slices.Collect].
func Collect[V any](seq iter.Seq[iter.Seq[V]]) iter.Seq[[]V] {
	return adapter.Map(slices.Collect, seq)
}

// Collect2 maps seq by [slices.Collect].
func Collect2[K, V any](seq iter.Seq2[iter.Seq[K], iter.Seq[V]]) iter.Seq2[[]K, []V] {
	return adapter.Map2(
		func(k iter.Seq[K], v iter.Seq[V]) ([]K, []V) {
			return slices.Collect(k), slices.Collect(v)
		},
		seq,
	)
}

// Clone maps seq by [slices.Clone].
func Clone[S ~[]E, E any](seq iter.Seq[S]) iter.Seq[S] {
	return adapter.Map(slices.Clone, seq)
}

// Clone2 maps seq by [slices.Clone].
func Clone2[S1 ~[]E1, S2 ~[]E2, E1, E2 any](seq iter.Seq2[S1, S2]) iter.Seq2[S1, S2] {
	return adapter.Map2(
		func(s1 S1, s2 S2) (S1, S2) { return slices.Clone(s1), slices.Clone(s2) },
		seq,
	)
}
