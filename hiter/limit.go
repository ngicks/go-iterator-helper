package hiter

import "iter"

// LimitUntil returns an iterator over seq that yields until f returns false.
func LimitUntil[V any](f func(V) bool, seq iter.Seq[V]) iter.Seq[V] {
	return func(yield func(V) bool) {
		for v := range seq {
			if !f(v) || !yield(v) {
				return
			}
		}
	}
}

// LimitUntil2 returns an iterator over seq that yields until f returns false.
func LimitUntil2[K, V any](f func(K, V) bool, seq iter.Seq2[K, V]) iter.Seq2[K, V] {
	return func(yield func(K, V) bool) {
		for k, v := range seq {
			if !f(k, v) || !yield(k, v) {
				return
			}
		}
	}
}

// LimitAfter is like [LimitUntil] but also yields the first value dissatisfying f(v).
func LimitAfter[V any](f func(V) bool, seq iter.Seq[V]) iter.Seq[V] {
	return func(yield func(V) bool) {
		for v := range seq {
			if !f(v) {
				yield(v)
				return
			}
			if !yield(v) {
				return
			}
		}
	}
}

// LimitAfter2 is like [LimitUntil2] but also yields the first pair dissatisfying f(k, v).
func LimitAfter2[K, V any](f func(K, V) bool, seq iter.Seq2[K, V]) iter.Seq2[K, V] {
	return func(yield func(K, V) bool) {
		for k, v := range seq {
			if !f(k, v) {
				yield(k, v)
				return
			}
			if !yield(k, v) {
				return
			}
		}
	}
}
