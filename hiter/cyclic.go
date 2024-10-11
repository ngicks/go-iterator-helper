package hiter

import "iter"

// Cycle converts seq into an infinite iterator by repeatedly calling seq.
// seq is assumed to be finite but pure and reusable.
// The iterator may yield forever if seq is repeatable; stopping it is caller's responsibility.
func Cycle[V any](seq iter.Seq[V]) iter.Seq[V] {
	return func(yield func(V) bool) {
		for {
			for v := range seq {
				if !yield(v) {
					return
				}
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
		for _, v := range buf {
			if !yield(v) {
				return
			}
		}
	}
}
