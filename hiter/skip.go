package hiter

import (
	"iter"
)

// SkipWhile returns an iterator over seq that skips n elements from seq.
func Skip[V any](seq iter.Seq[V], n int) iter.Seq[V] {
	return func(yield func(V) bool) {
		if n <= 0 {
			return
		}
		for v := range seq {
			if n--; n >= 0 {
				continue
			}
			if !yield(v) {
				return
			}
		}
	}
}

// SkipWhile returns an iterator over seq that skips n key-value pairs from seq.
func Skip2[K, V any](seq iter.Seq2[K, V], n int) iter.Seq2[K, V] {
	return func(yield func(K, V) bool) {
		if n <= 0 {
			return
		}
		for k, v := range seq {
			if n--; n >= 0 {
				continue
			}
			if !yield(k, v) {
				return
			}
		}
	}
}

// SkipWhile returns an iterator over seq that skips elements until f returns false.
func SkipWhile[V any](seq iter.Seq[V], f func(V) bool) iter.Seq[V] {
	return func(yield func(V) bool) {
		skipping := true
		for v := range seq {
			if skipping && !f(v) {
				continue
			}
			skipping = false
			if !yield(v) {
				return
			}
		}
	}
}

// SkipWhile2 returns an iterator over seq that skips key-value pairs until f returns false.
func SkipWhile2[K, V any](seq iter.Seq2[K, V], f func(K, V) bool) iter.Seq2[K, V] {
	return func(yield func(K, V) bool) {
		skipping := true
		for k, v := range seq {
			if skipping && !f(k, v) {
				continue
			}
			skipping = false
			if !yield(k, v) {
				return
			}
		}
	}
}
