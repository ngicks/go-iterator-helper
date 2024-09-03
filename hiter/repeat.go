package hiter

import "iter"

// Repeat returns an iterator that generates v n times.
func Repeat[V any](v V, n int) iter.Seq[V] {
	return func(yield func(V) bool) {
		for range n {
			if !yield(v) {
				return
			}
		}
	}
}

// Repeat2 returns an iterator that generates v n times.
func Repeat2[K, V any](k K, v V, n int) iter.Seq2[K, V] {
	return func(yield func(K, V) bool) {
		for range n {
			if !yield(k, v) {
				return
			}
		}
	}
}
