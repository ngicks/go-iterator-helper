package tee

import "iter"

// Experimental: not tested and might be changed any time.
func Push[V any](seq iter.Seq[V], f ...func(V) bool) bool {
	for v := range seq {
		for _, f := range f {
			if !f(v) {
				return false
			}
		}
	}
	return true
}

// Experimental: not tested and might be changed any time.
func Push2[K, V any](seq iter.Seq2[K, V], f ...func(K, V) bool) bool {
	for k, v := range seq {
		for _, f := range f {
			if !f(k, v) {
				return false
			}
		}
	}
	return true
}
