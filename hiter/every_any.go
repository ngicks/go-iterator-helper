package hiter

import "iter"

// Every checks every values from seq satisfying fn.
// Every return false immediately after it found the value dissatisfying fn,
// otherwise returns true.
func Every[V any](fn func(V) bool, seq iter.Seq[V]) bool {
	for v := range seq {
		if !fn(v) {
			return false
		}
	}
	return true
}

// Every2 checks every key-value pairs from seq satisfying fn.
// Every2 return false immediately after it found the pair dissatisfying fn,
// otherwise returns true.
func Every2[K, V any](fn func(K, V) bool, seq iter.Seq2[K, V]) bool {
	for k, v := range seq {
		if !fn(k, v) {
			return false
		}
	}
	return true
}

// Any returns true if any of values from seq satisfies fn.
// Otherwise it returns false.
func Any[V any](fn func(V) bool, seq iter.Seq[V]) bool {
	_, idx := FindFunc(fn, seq)
	return idx >= 0
}

// Any2 returns true if any of key-value pairs from seq satisfies fn.
// Otherwise it returns false.
func Any2[K, V any](fn func(K, V) bool, seq iter.Seq2[K, V]) bool {
	_, _, idx := FindFunc2(fn, seq)
	return idx >= 0
}
