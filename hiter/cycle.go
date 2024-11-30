package hiter

import "iter"

// Cycle converts seq into an infinite iterator by repeatedly calling seq.
// seq is assumed to be finite but pure and reusable.
// The iterator may yield forever if seq is repeatable; stopping it is caller's responsibility.
func Cycle[V any](seq iter.Seq[V]) iter.Seq[V] {
	return func(yield func(V) bool) {
		for {
			atLeastOne := false
			for v := range seq {
				atLeastOne = true
				if !yield(v) {
					return
				}
			}
			if !atLeastOne {
				break
			}
		}
	}
}

// Cycle2 converts seq into an infinite iterator by repeatedly calling seq.
// seq is assumed to be finite but pure and reusable.
// The iterator may yield forever if seq is repeatable; stopping it is caller's responsibility.
func Cycle2[K, V any](seq iter.Seq2[K, V]) iter.Seq2[K, V] {
	return func(yield func(K, V) bool) {
		for {
			atLeastOne := false
			for k, v := range seq {
				atLeastOne = true
				if !yield(k, v) {
					return
				}
			}
			if !atLeastOne {
				break
			}
		}
	}
}

// CycleBuffered is like [Cycle] but seq is called only once.
// Values from seq is buffered and from second time and on,
// the iterator uses buffered contents.
//
// seq must be finite and small, otherwise huge amount of memory will be consumed.
func CycleBuffered[V any](seq iter.Seq[V]) iter.Seq[V] {
	return func(yield func(V) bool) {
		var buf []V
		for v := range seq {
			if !yield(v) {
				return
			}
			buf = append(buf, v)
		}
		if len(buf) == 0 {
			return
		}
		for {
			for _, v := range buf {
				if !yield(v) {
					return
				}
			}
		}
	}
}

// CycleBuffered2 is like [Cycle2] but seq is called only once.
// Key-Value pairs from seq is buffered and from second time and on,
// the iterator uses buffered contents.
//
// seq must be finite and small, otherwise huge amount of memory will be consumed.
func CycleBuffered2[K, V any](seq iter.Seq2[K, V]) iter.Seq2[K, V] {
	return func(yield func(K, V) bool) {
		var buf []KeyValue[K, V]
		for k, v := range seq {
			if !yield(k, v) {
				return
			}
			buf = append(buf, KeyValue[K, V]{k, v})
		}
		if len(buf) == 0 {
			return
		}
		for {
			for _, kv := range buf {
				if !yield(kv.K, kv.V) {
					return
				}
			}
		}
	}
}
