package hiter

import "iter"

func Tap[V any](tap func(V), seq iter.Seq[V]) iter.Seq[V] {
	return func(yield func(V) bool) {
		for v := range seq {
			tap(v)
			if !yield(v) {
				return
			}
		}
	}
}

func Tap2[K, V any](tap func(K, V), seq iter.Seq2[K, V]) iter.Seq2[K, V] {
	return func(yield func(K, V) bool) {
		for k, v := range seq {
			tap(k, v)
			if !yield(k, v) {
				return
			}
		}
	}
}
