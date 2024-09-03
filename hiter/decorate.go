package hiter

import "iter"

// Decorate decorates seq by prepend and append,
// by yielding additional elements before and after seq yields.
func Decorate[V any](seq iter.Seq[V], prepend, append Iterable[V]) iter.Seq[V] {
	return func(yield func(V) bool) {
		for v := range seq {
			if prepend != nil {
				for vp := range prepend.Iter() {
					if !yield(vp) {
						return
					}
				}
			}
			if !yield(v) {
				return
			}
			if append != nil {
				for va := range append.Iter() {
					if !yield(va) {
						return
					}
				}
			}
		}
	}
}

// Decorate2 decorates seq by prepend and append,
// by yielding additional elements before and after seq yields.
func Decorate2[K, V any](seq iter.Seq2[K, V], prepend, append Iterable2[K, V]) iter.Seq2[K, V] {
	return func(yield func(K, V) bool) {
		for k, v := range seq {
			if prepend != nil {
				for kp, vp := range prepend.Iter2() {
					if !yield(kp, vp) {
						return
					}
				}
			}
			if !yield(k, v) {
				return
			}
			if append != nil {
				for ka, va := range append.Iter2() {
					if !yield(ka, va) {
						return
					}
				}
			}
		}
	}
}
