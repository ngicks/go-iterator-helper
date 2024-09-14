package hiter

import (
	"iter"
)

// Window returns an iterator over overlapping sub-slices of n size (moving windows).
// n must be a positive non zero value.
// Values from the iterator are always size of n.
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
