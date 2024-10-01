package hiter

import "iter"

func Every[V any](fn func(V) bool, seq iter.Seq[V]) bool {
	for v := range seq {
		if !fn(v) {
			return false
		}
	}
	return true
}

func Every2[K, V any](fn func(K, V) bool, seq iter.Seq2[K, V]) bool {
	for k, v := range seq {
		if !fn(k, v) {
			return false
		}
	}
	return true
}

func Any[V any](fn func(V) bool, seq iter.Seq[V]) bool {
	_, idx := FindFunc(fn, seq)
	return idx >= 0
}

func Any2[K, V any](fn func(K, V) bool, seq iter.Seq2[K, V]) bool {
	_, _, idx := FindFunc2(fn, seq)
	return idx >= 0
}
