package hiter

import "iter"

func Single[V any](v V) iter.Seq[V] {
	return func(yield func(V) bool) {
		yield(v)
	}
}

func Single2[K, V any](k K, v V) iter.Seq2[K, V] {
	return func(yield func(K, V) bool) {
		yield(k, v)
	}
}
