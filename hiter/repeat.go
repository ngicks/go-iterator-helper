package hiter

import "iter"

// Repeat returns an iterator over v repeated n times.
// If n < 0, the returned iterator repeats forever.
func Repeat[V any](v V, n int) iter.Seq[V] {
	if n < 0 {
		return func(yield func(V) bool) {
			for {
				if !yield(v) {
					return
				}
			}
		}
	}
	return func(yield func(V) bool) {
		// no state in the seq.
		for n := n; n != 0; n-- {
			if !yield(v) {
				return
			}
		}
	}
}

// Repeat2 returns an iterator over the pair of k and v repeated n times.
// If n < 0, the returned iterator repeats forever.
func Repeat2[K, V any](k K, v V, n int) iter.Seq2[K, V] {
	if n < 0 {
		return func(yield func(K, V) bool) {
			for {
				if !yield(k, v) {
					return
				}
			}
		}
	}
	return func(yield func(K, V) bool) {
		// no state in the seq.
		for n := n; n != 0; n-- {
			if !yield(k, v) {
				return
			}
		}
	}
}

// RepeatFunc returns an iterator that generates result from fnV n times.
// If n < 0, the returned iterator repeats forever.
func RepeatFunc[V any](fnV func() V, n int) iter.Seq[V] {
	if n < 0 {
		return func(yield func(V) bool) {
			for {
				if !yield(fnV()) {
					return
				}
			}
		}
	}
	return func(yield func(V) bool) {
		for n := n; n != 0; n-- {
			if !yield(fnV()) {
				return
			}
		}
	}
}

// RepeatFunc2 returns an iterator that generates result of fnK and fnV n times.
// If n < 0, the returned iterator repeats forever.
func RepeatFunc2[K, V any](fnK func() K, fnV func() V, n int) iter.Seq2[K, V] {
	if n < 0 {
		return func(yield func(K, V) bool) {
			for {
				if !yield(fnK(), fnV()) {
					return
				}
			}
		}
	}
	return func(yield func(K, V) bool) {
		for n := n; n != 0; n-- {
			if !yield(fnK(), fnV()) {
				return
			}
		}
	}
}
