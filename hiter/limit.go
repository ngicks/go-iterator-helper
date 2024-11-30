package hiter

import "iter"

// Limit returns an iterator over seq that stops after n values.
//
// Note: Limit is redefined since xiter's implementation is stateful.
func Limit[V any](n int, seq iter.Seq[V]) iter.Seq[V] {
	return func(yield func(V) bool) {
		n := n // no state in the iterator
		if n <= 0 {
			return
		}
		for v := range seq {
			if !yield(v) {
				return
			}
			if n--; n <= 0 {
				break
			}
		}
	}
}

// Limit2 returns an iterator over seq that stops after n key-value pairs.
//
// Note: Limit2 is redefined since xiter's implementation is stateful.
func Limit2[K, V any](n int, seq iter.Seq2[K, V]) iter.Seq2[K, V] {
	return func(yield func(K, V) bool) {
		n := n
		if n <= 0 {
			return
		}
		for k, v := range seq {
			if !yield(k, v) {
				return
			}
			if n--; n <= 0 {
				break
			}
		}
	}
}

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
