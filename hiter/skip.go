package hiter

import (
	"iter"
)

// SkipWhile returns an iterator over seq that skips n elements from seq.
func Skip[V any](n int, seq iter.Seq[V]) iter.Seq[V] {
	return func(yield func(V) bool) {
		if n <= 0 {
			return
		}
		for v := range seq {
			if n--; n >= 0 {
				continue
			}
			if !yield(v) {
				return
			}
		}
	}
}

// SkipWhile returns an iterator over seq that skips n key-value pairs from seq.
func Skip2[K, V any](n int, seq iter.Seq2[K, V]) iter.Seq2[K, V] {
	return func(yield func(K, V) bool) {
		if n <= 0 {
			return
		}
		for k, v := range seq {
			if n--; n >= 0 {
				continue
			}
			if !yield(k, v) {
				return
			}
		}
	}
}

// SkipLast returns an iterator over seq that skips last n elements.
func SkipLast[V any](n int, seq iter.Seq[V]) iter.Seq[V] {
	return func(yield func(V) bool) {
		var ( // easy implementation for ring buffer.
			buf    = make([]V, n)
			cursor int
			full   bool
		)
		for v := range seq {
			if !full {
				buf[cursor] = v
				cursor++
				if cursor == n {
					cursor = 0
					full = true
				}
				continue
			}
			vOld := buf[cursor]
			if !yield(vOld) {
				return
			}
			buf[cursor] = v
			cursor++
			if cursor == n {
				cursor = 0
			}
		}
	}
}

// SkipLast returns an iterator over seq that skips last n key-value pairs.
func SkipLast2[K, V any](n int, seq iter.Seq2[K, V]) iter.Seq2[K, V] {
	return func(yield func(K, V) bool) {
		var ( // easy implementation for ring buffer.
			buf    = make([]KeyValue[K, V], n)
			cursor int
			full   bool
		)
		for k, v := range seq {
			if !full {
				buf[cursor] = KeyValue[K, V]{k, v}
				cursor++
				if cursor == n {
					cursor = 0
					full = true
				}
				continue
			}
			kvOld := buf[cursor]
			if !yield(kvOld.K, kvOld.V) {
				return
			}
			buf[cursor] = KeyValue[K, V]{k, v}
			cursor++
			if cursor == n {
				cursor = 0
			}
		}
	}
}

// SkipWhile returns an iterator over seq that skips elements until f returns false.
func SkipWhile[V any](f func(V) bool, seq iter.Seq[V]) iter.Seq[V] {
	return func(yield func(V) bool) {
		skipping := true
		for v := range seq {
			if skipping && !f(v) {
				continue
			}
			skipping = false
			if !yield(v) {
				return
			}
		}
	}
}

// SkipWhile2 returns an iterator over seq that skips key-value pairs until f returns false.
func SkipWhile2[K, V any](f func(K, V) bool, seq iter.Seq2[K, V]) iter.Seq2[K, V] {
	return func(yield func(K, V) bool) {
		skipping := true
		for k, v := range seq {
			if skipping && !f(k, v) {
				continue
			}
			skipping = false
			if !yield(k, v) {
				return
			}
		}
	}
}
