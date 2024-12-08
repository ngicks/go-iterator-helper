package tee

import (
	"iter"
)

// MultiPusher is [iter.Seq] equivalent of [io.MultiWriter].
//
// Experimental: not tested and might be changed any time.
func MultiPusher[V any](pushers ...func(V) bool) func(v V) bool {
	return func(v V) bool {
		for _, p := range pushers {
			if !p(v) {
				return false
			}
		}
		return true
	}
}

// MultiPusher2 is [iter.Seq2] equivalent of [io.MultiWriter].
//
// Experimental: not tested and might be changed any time.
func MultiPusher2[K, V any](pushers ...func(K, V) bool) func(k K, v V) bool {
	return func(k K, v V) bool {
		for _, p := range pushers {
			if !p(k, v) {
				return false
			}
		}
		return true
	}
}

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
