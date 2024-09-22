package hiter

import (
	"iter"
)

// Window returns an iterator over overlapping sub-slices of n size (moving windows).
// n must be a positive non zero value.
// Values from the iterator are always slices of n size.
// The iterator yields nothing when it is not possible.
func Window[S ~[]E, E any](s S, n int) iter.Seq[S] {
	return func(yield func(S) bool) {
		if n <= 0 {
			return
		}
		var (
			start = 0
			end   = n
		)
		for {
			if end > len(s) {
				return
			}
			if !yield(s[start:end]) {
				return
			}
			start++
			end++
		}
	}
}

// WindowSeq buffers n elements from seq, then yields an iterator over them.
// The iterator first fills the buffer with values from seq in FIFO manner.
// After full, it yields another iterator over buffered values.
// Each time value is yielded from seq, the iterator pushes value to tail of the buffer
// and yields an iterator over the updated buffer.
//
// n must be a positive non zero value.
// Each iterator yields exact n size of values.
// If seq yields less than n, the iterator yields nothing.
func WindowSeq[V any](n int, seq iter.Seq[V]) iter.Seq[iter.Seq[V]] {
	return func(yield func(iter.Seq[V]) bool) {
		if n <= 0 {
			return
		}
		var (
			buf    = make([]V, n)
			cursor = 0
			full   = false
		)
		for e := range seq {
			if !full {
				buf[cursor] = e
				cursor++
				if cursor == n {
					cursor = 0
					full = true
					if !yield(sliceRing(buf, cursor)) {
						return
					}
				}
				continue
			}
			buf[cursor] = e
			cursor++
			if cursor == n {
				cursor = 0
			}
			if !yield(sliceRing(buf, cursor)) {
				return
			}
		}
	}
}

func sliceRing[S ~[]E, E any](s S, start int) iter.Seq[E] {
	return func(yield func(E) bool) {
		if !yield(s[start]) {
			return
		}
		for i := start + 1; ; i++ {
			if i >= len(s) {
				i = 0
			}
			if i == start {
				break
			}
			if !yield(s[i]) {
				return
			}
		}
	}
}
