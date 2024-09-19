package hiter

import (
	"iter"
)

// Compact skips consecutive runs of equal elements from seq.
// The returned iterator is pure and stateless as long as seq is so.
func Compact[V comparable](seq iter.Seq[V]) iter.Seq[V] {
	return func(yield func(V) bool) {
		var (
			first bool = true
			prev  V
		)
		for t := range seq {
			if first {
				first = false
				if !yield(t) {
					return
				}
			} else if prev != t {
				if !yield(t) {
					return
				}
			}
			prev = t
		}
	}
}

// CompactFunc is like [Compact] but uses an equality function to compare elements.
// For runs of elements that compare equal, CompactFunc keeps the first one.
func CompactFunc[V any](eq func(i, j V) bool, seq iter.Seq[V]) iter.Seq[V] {
	return func(yield func(V) bool) {
		var (
			first bool = true
			prev  V
		)
		for t := range seq {
			if first {
				first = false
				if !yield(t) {
					return
				}
			} else if !eq(prev, t) {
				if !yield(t) {
					return
				}
			}
			prev = t
		}
	}
}

// Compact2 skips consecutive runs of equal k-v pairs from seq.
// The returned iterator is pure and stateless as long as seq is so.
func Compact2[K, V comparable](seq iter.Seq2[K, V]) iter.Seq2[K, V] {
	return func(yield func(K, V) bool) {
		var (
			first bool = true
			prevK K
			prevV V
		)
		for k, v := range seq {
			if first {
				first = false
				if !yield(k, v) {
					return
				}
			} else if prevK != k || prevV != v {
				if !yield(k, v) {
					return
				}
			}
			prevK = k
			prevV = v
		}
	}
}

// CompactFunc2 is like [Compact2] but uses an equality function to compare elements.
// For runs of elements that compare equal, CompactFunc2 keeps the first one.
func CompactFunc2[K, V any](eq func(k1 K, v1 V, k2 K, v2 V) bool, seq iter.Seq2[K, V]) iter.Seq2[K, V] {
	return func(yield func(K, V) bool) {
		var (
			first bool = true
			prevK K
			prevV V
		)
		for k, v := range seq {
			if first {
				first = false
				if !yield(k, v) {
					return
				}
			} else if !eq(prevK, prevV, k, v) {
				if !yield(k, v) {
					return
				}
			}
			prevK = k
			prevV = v
		}
	}
}
