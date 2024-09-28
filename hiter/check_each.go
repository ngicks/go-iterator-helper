package hiter

import "iter"

// CheckEach returns an iterator over seq.
// It also calls check after each n values yielded from seq.
// fn receives a value from seq and the current count of it to inspect and decide to stop.
// n is capped to 1 if it is less.
func CheckEach[V any](n int, check func(v V, i int) bool, seq iter.Seq[V]) iter.Seq[V] {
	if n <= 1 {
		// The specialized case...Does it have any effect?
		return func(yield func(V) bool) {
			i := 0
			for v := range seq {
				if !check(v, i) {
					return
				}
				i++
				if !yield(v) {
					return
				}
			}
		}
	}
	return func(yield func(V) bool) {
		i := 0
		nn := n
		for v := range seq {
			nn--
			if nn == 0 {
				if !check(v, i) {
					return
				}
				nn = n
			}
			i++
			if !yield(v) {
				return
			}
		}
	}
}

// CheckEach2 returns an iterator over seq.
// It also calls check after each n key-value pairs yielded from seq.
// fn receives a key-value pair from seq and the current count of it to inspect and decide to stop.
// n is capped to 1 if it is less.
func CheckEach2[K, V any](n int, check func(k K, v V, i int) bool, seq iter.Seq2[K, V]) iter.Seq2[K, V] {
	if n <= 1 {
		return func(yield func(K, V) bool) {
			i := 0
			for k, v := range seq {
				if !check(k, v, i) {
					return
				}
				i++
				if !yield(k, v) {
					return
				}
			}
		}
	}
	return func(yield func(K, V) bool) {
		i := 0
		nn := n
		for k, v := range seq {
			nn--
			if nn == 0 {
				if !check(k, v, i) {
					return
				}
				nn = n
			}
			i++
			if !yield(k, v) {
				return
			}
		}
	}
}
