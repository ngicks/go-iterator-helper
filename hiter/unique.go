package hiter

import (
	"iter"
)

// UniqueFunc returns an iterator over every unique values from seq.
// f returns unique identifiers for v.
//
// Unlike [Compact], UniqueFunc buffers ids.
// Therefore the larger cardinality, the more memory is consumed.
func UniqueFunc[V any, I comparable](f func(V) I, seq iter.Seq[V]) iter.Seq[V] {
	return func(yield func(V) bool) {
		seen := map[I]bool{}
		for v := range seq {
			i := f(v)
			if seen[i] {
				continue
			}
			seen[i] = true
			if !yield(v) {
				return
			}
		}
	}
}

// Unique is [UniqueFunc] where the value itself is an identifier.
func Unique[V comparable](seq iter.Seq[V]) iter.Seq[V] {
	return UniqueFunc(func(v V) V { return v }, seq)
}

// UniqueFunc2 returns an iterator over every unique key-value pairs from seq.
// f returns unique identifiers for combination of k and v.
//
// Unlike [Compact2], UniqueFunc2 buffers ids.
// Therefore the larger cardinality, the more memory is consumed.
func UniqueFunc2[K, V any, I comparable](f func(K, V) I, seq iter.Seq2[K, V]) iter.Seq2[K, V] {
	return func(yield func(K, V) bool) {
		seen := map[I]bool{}
		for k, v := range seq {
			i := f(k, v)
			if seen[i] {
				continue
			}
			seen[i] = true
			if !yield(k, v) {
				return
			}
		}
	}
}

// Unique2 is [UniqueFunc2] where the key-value pair itself is an identifier.
func Unique2[K, V comparable](seq iter.Seq2[K, V]) iter.Seq2[K, V] {
	return UniqueFunc2(func(k K, v V) KeyValue[K, V] { return KeyValue[K, V]{k, v} }, seq)
}
