package hiter

import "iter"

// Tap returns an iterator over seq without any modification to values from seq.
// tap is called against every value from seq to observe values.
// For purpose of Tap, the tap callback function should not retain arguments.
func Tap[V any](tap func(V), seq iter.Seq[V]) iter.Seq[V] {
	return func(yield func(V) bool) {
		for v := range seq {
			tap(v)
			if !yield(v) {
				return
			}
		}
	}
}

// Tap2 returns an iterator over seq without any modification to pairs from seq.
// tap is called against every pair from seq to observe pairs.
// For purpose of Tap, the tap callback function should not retain arguments.
func Tap2[K, V any](tap func(K, V), seq iter.Seq2[K, V]) iter.Seq2[K, V] {
	return func(yield func(K, V) bool) {
		for k, v := range seq {
			tap(k, v)
			if !yield(k, v) {
				return
			}
		}
	}
}

// TapLast returns an iterator over values from seq.
// tap is executed after all values from seq are consumed without breaking it.
//
// Experimental: not tested and might be changed any time.
func TapLast[V any](tap func(), seq iter.Seq[V]) iter.Seq[V] {
	return func(yield func(V) bool) {
		for v := range seq {
			if !yield(v) {
				return
			}
		}
		tap()
	}
}

// TapLast2 returns an iterator over key-value pairs from seq.
// tap is executed after all pairs from seq are consumed without breaking it.
//
// Experimental: not tested and might be changed any time.
func TapLast2[K, V any](tap func(), seq iter.Seq2[K, V]) iter.Seq2[K, V] {
	return func(yield func(K, V) bool) {
		for k, v := range seq {
			if !yield(k, v) {
				return
			}
		}
		tap()
	}
}
