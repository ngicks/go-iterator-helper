package hiter

import (
	"iter"
)

type KeyValue[K, V any] struct {
	K K
	V V
}

// AppendSeq2 appends the values from seq to the KeyValue slice and
// returns the extended slice.
func AppendSeq2[S ~[]KeyValue[K, V], K, V any](s S, seq iter.Seq2[K, V]) S {
	for k, v := range seq {
		s = append(s, KeyValue[K, V]{k, v})
	}
	return s
}

// Collect2 collects values from seq into a new KeyValue slice and returns it.
func Collect2[K, V any](seq iter.Seq2[K, V]) []KeyValue[K, V] {
	return AppendSeq2[[]KeyValue[K, V]](nil, seq)
}

// ToKeyValue converts [iter.Seq2][K, V] into iter.Seq[KeyValue[K, V]].
// This functions is particularly useful when sending values from [iter.Seq2][K, V] through
// some data transfer mechanism that only allows data to be single value, e.g. channels.
func ToKeyValue[K, V any](seq iter.Seq2[K, V]) iter.Seq[KeyValue[K, V]] {
	return func(yield func(KeyValue[K, V]) bool) {
		for k, v := range seq {
			if !yield(KeyValue[K, V]{k, v}) {
				return
			}
		}
	}
}

// FromKeyValue unwraps iter.Seq[KeyValue[K, V]] into iter.Seq2[K, V] to counter-part,
// serving a counterpart for [ToKeyValue].
//
// In case values from seq needs to be sent through some data transfer mechanism
// that only allows data to be single value, like channels,
// some caller might decide to wrap values into KeyValue[K, V], maybe by [ToKeyValue].
// If target helpers only accept iter.Seq2[K, V], then FromKeyValues is useful.
func FromKeyValue[K, V any](seq iter.Seq[KeyValue[K, V]]) iter.Seq2[K, V] {
	return func(yield func(K, V) bool) {
		for kv := range seq {
			if !yield(kv.K, kv.V) {
				return
			}
		}
	}
}
