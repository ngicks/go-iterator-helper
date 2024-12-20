package hiter

import (
	"iter"
)

// Step returns an iterator over numerics values starting from initial and added step at each step.
// The iterator iterates forever. The caller might want to limit it by [Limit].
func Step[N Numeric](initial, step N) iter.Seq[N] {
	return func(yield func(N) bool) {
		for n := initial; ; n += step {
			if !yield(n) {
				return
			}
		}
	}
}

// StepBy returns an iterator over pair of index and value associated the index.
// The index starts from initial and steps by step.
func StepBy[V any](initial, step int, v []V) iter.Seq2[int, V] {
	return func(yield func(int, V) bool) {
		if initial < 0 {
			return
		}
		for i := initial; 0 <= i && i < len(v); i += step {
			if !yield(i, v[i]) {
				return
			}
		}
	}
}
