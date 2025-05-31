package hiter

import "iter"

// Enumerate adds 0-indexed integer counts to every values from seq.
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

// Pairs combines x and y into an iterator over key-value pairs.
// Unlike [Zip], if either stops, the returned iterator stops.
func Pairs[K, V any](x iter.Seq[K], y iter.Seq[V]) iter.Seq2[K, V] {
	return func(yield func(K, V) bool) {
		next, stop := iter.Pull(y)
		defer stop()
		for k := range x {
			v, ok := next()
			if !ok || !yield(k, v) {
				return
			}
		}
	}
}

// Pairs2 combines x and y into an iterator over pairs of key-value pairs.
// Unlike [Zip2], if either stops, the returned iterator stops.
func Pairs2[K1, V1, K2, V2 any](x iter.Seq2[K1, V1], y iter.Seq2[K2, V2]) iter.Seq2[KeyValue[K1, V1], KeyValue[K2, V2]] {
	return func(yield func(KeyValue[K1, V1], KeyValue[K2, V2]) bool) {
		next, stop := iter.Pull2(y)
		defer stop()
		for k1, v1 := range x {
			k2, v2, ok := next()
			if !ok || !yield(KVPair(k1, v1), KVPair(k2, v2)) {
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
func OmitL[K, V any](seq iter.Seq2[K, V]) iter.Seq[K] {
	return func(yield func(K) bool) {
		for t := range seq {
			if !yield(t) {
				return
			}
		}
	}
}

// OmitF drops former part of key-value pairs that seq generates.
func OmitF[K, V any](seq iter.Seq2[K, V]) iter.Seq[V] {
	return func(yield func(V) bool) {
		for _, u := range seq {
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

// Unify unifies key-value pairs from seq into single values by applying fn
// to convert [iter.Seq2][K, V] to [iter.Seq][U].
func Unify[K, V, U any](fn func(K, V) U, seq iter.Seq2[K, V]) iter.Seq[U] {
	return func(yield func(U) bool) {
		for k, v := range seq {
			if !yield(fn(k, v)) {
				return
			}
		}
	}
}

// Divide splits values from seq to key-value pairs by applying fn
// to convert [iter.Seq][U] to [iter.Seq2][K, V].
func Divide[K, V, U any](fn func(U) (K, V), seq iter.Seq[U]) iter.Seq2[K, V] {
	return func(yield func(K, V) bool) {
		for u := range seq {
			if !yield(fn(u)) {
				return
			}
		}
	}
}
