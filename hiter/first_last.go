package hiter

import "iter"

// First returns the first value from seq.
// ok is false if seq yields no value.
func First[V any](seq iter.Seq[V]) (k V, ok bool) {
	for v := range seq {
		return v, true
	}
	return *new(V), false
}

// First2 returns the first key-value pair from seq.
// ok is false if seq yields no pair.
func First2[K, V any](seq iter.Seq2[K, V]) (k K, v V, ok bool) {
	for k, v := range seq {
		return k, v, true
	}
	return *new(K), *new(V), false
}

// Last returns the last value from seq.
// ok is false if seq yields no value.
func Last[V any](seq iter.Seq[V]) (v V, ok bool) {
	for v = range seq {
		ok = true
	}
	return
}

// Last2 returns the last key-value pair from seq.
// ok is false if seq yields no pair.
func Last2[K, V any](seq iter.Seq2[K, V]) (k K, v V, ok bool) {
	for k, v = range seq {
		ok = true
	}
	return
}
