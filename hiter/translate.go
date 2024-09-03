package hiter

import "iter"

// Enumerate wraps seq so that former part of paired values has an index starting from 0.
// Each time values are yielded index is increased by 1.
func Enumerate[T any](seq iter.Seq[T]) iter.Seq2[int, T] {
	return func(yield func(int, T) bool) {
		idx := 0
		for t := range seq {
			if !yield(idx, t) {
				return
			}
			idx++
		}
	}
}

// Combine combines seq1 and seq2 into single key-value pairs.
func Combine[K, V any](seq1 iter.Seq[K], seq2 iter.Seq[V]) iter.Seq2[K, V] {
	return func(yield func(K, V) bool) {
		next, stop := iter.Pull(seq2)
		defer stop()
		for k := range seq1 {
			v, ok := next()
			if !ok || !yield(k, v) {
				return
			}
		}
	}
}

// Transpose returns an iterator over seq that yields K and V reversed.
func Transpose[K, V any](seq iter.Seq2[K, V]) iter.Seq2[V, K] {
	return func(yield func(V, K) bool) {
		for t, u := range seq {
			if !yield(u, t) {
				return
			}
		}
	}
}

// OmitL drops latter part of key-value pairs that seq generates.
func OmitL[T, U any](i iter.Seq2[T, U]) iter.Seq[T] {
	return func(yield func(T) bool) {
		for t := range i {
			if !yield(t) {
				return
			}
		}
	}
}

// OmitF drops former part of key-value pairs that seq generates.
func OmitF[T, U any](i iter.Seq2[T, U]) iter.Seq[U] {
	return func(yield func(U) bool) {
		for _, u := range i {
			if !yield(u) {
				return
			}
		}
	}
}

// Omit returns an iterator over seq but drops data seq yields.
func Omit[K any](seq iter.Seq[K]) func(yield func() bool) {
	return func(yield func() bool) {
		for range seq {
			if !yield() {
				return
			}
		}
	}
}

// Omit2 returns an iterator over seq but drops data seq yields.
func Omit2[K, V any](seq iter.Seq2[K, V]) func(yield func() bool) {
	return func(yield func() bool) {
		for range seq {
			if !yield() {
				return
			}
		}
	}
}
