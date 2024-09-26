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
			if !yield(s[start:end:end]) {
				return
			}
			start++
			end++
		}
	}
}

// WindowSeq allocates n sized buffer and fills it with values from seq in FIFO-manner.
// Once the buffer is full, the returned iterator yields iterator over buffered values
// each time value is yielded from the input seq.
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
			cursor = (cursor + 1) % n
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
			i = i % len(s)
			if i == start {
				break
			}
			if !yield(s[i]) {
				return
			}
		}
	}
}
