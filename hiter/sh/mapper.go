// package sh defines some short hand iterator adapters.
// sh only holds functions which only combines other element in this module.
package sh

import (
	"iter"
	"reflect"
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

func Clone[S ~[]E, E any](seq iter.Seq[S]) iter.Seq[S] {
	return mapIter(slices.Clone, seq)
}

