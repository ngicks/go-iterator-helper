package hiter

import (
	"iter"
)

// Say we have sequence data of n size.
// |++++++++++data++++++++++|
// Limit limits to n elements from head.
// |+++++++:n---------------|
// Skip skips n elements from head
// |--------:n++++++++++++++|
// SkipLast skips n element from tails
// |++++++++++++++n:--------|
//
// Soooo... will we have LimitLast?
// It should skip element to ensure only size of n elements from *tail* are returned to caller.
// |---------------n:+++++++|
// Definitely what I should refer to is that,
// it may block so long time.
// Normally we should not have control how many elements the seq is supposed to yield.
// Implementation should look like
//
// func LimitLast[V any](n int, seq iter.Seq[V]) iter.Seq[V] {
// 	return func(yield func(V) bool) {
// 		var (
// 			buf    = make([]V, n)
// 			cursor int
// 			full   bool
// 		)
// 		for v := range seq {
// 			buf[cursor] = v
// 			cursor++
// 			if cursor == n {
// 				cursor = 0
// 				full = true
// 			}
// 		}
// 		end := cursor
// 		if full {
// 			end = n
// 		}
// 		for i := 0; i < end; i++ {
// 			if !yield(buf[i]) {
// 				return
// 			}
// 		}
// 	}
// }
//
// The uncontrolled long blocking sounds not too good to me.

// Skip returns an iterator over seq that skips n elements from seq.
func Skip[V any](n int, seq iter.Seq[V]) iter.Seq[V] {
	return func(yield func(V) bool) {
		n := n
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

// Skip2 returns an iterator over seq that skips n key-value pairs from seq.
func Skip2[K, V any](n int, seq iter.Seq2[K, V]) iter.Seq2[K, V] {
	return func(yield func(K, V) bool) {
		n := n // no state
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

// SkipLast returns an iterator over seq that skips last n elements of seq.
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

// SkipLast2 returns an iterator over seq that skips last n key-value pairs of seq.
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
			if skipping && f(v) {
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
			if skipping && f(k, v) {
				continue
			}
			skipping = false
			if !yield(k, v) {
				return
			}
		}
	}
}
