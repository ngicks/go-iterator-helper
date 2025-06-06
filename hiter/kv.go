package hiter

import "iter"

// KeyValue is combination of K and V.
// It is useful if paris of k and v needs to be sent through
// somewhere it only accepts a single vlaue, e.g. channel.
//
// There are some utilities that help you consume KeyValue easier:
//
//   - [KVPair] : creates a KeyValue which saves you some key strokes by inferring types.
//   - [Values2] : converts []KeyValue[K, V] based types to iter.Seq2.
//   - [AppendSeq2], [Collect2] : collects iter.Seq2 to []KeyValue or types based on it.
//   - [ToKeyValue], [FromKeyValue] : converts iter.Seq2[K, V] to/from iter.Seq[KeyValue[K, V]].
type KeyValue[K, V any] struct {
	K K
	V V
}

func (k KeyValue[K, V]) Unpack() (K, V) {
	return k.K, k.V
}

func KVPair[K, V any](k K, v V) KeyValue[K, V] {
	return KeyValue[K, V]{k, v}
}

// Values2 returns an iterator that yields the KeyValue slice elements in order.
func Values2[S ~[]KeyValue[K, V], K, V any](s S) iter.Seq2[K, V] {
	return func(yield func(K, V) bool) {
		for _, kv := range s {
			if !yield(kv.K, kv.V) {
				return
			}
		}
	}
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

// FromKeyValue unwraps iter.Seq[KeyValue[K, V]] into iter.Seq2[K, V]
// serving as a counterpart for [ToKeyValue].
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

var _ Iterable2[any, any] = KeyValues[any, any](nil)

// KeyValues adds the Iter2 method to slice of KeyValue-s.
type KeyValues[K, V any] []KeyValue[K, V]

func (v KeyValues[K, V]) Iter2() iter.Seq2[K, V] {
	return Values2(v)
}
