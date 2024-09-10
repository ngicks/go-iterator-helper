package hiter

import (
	"iter"
)

// Window returns an iterator over overlapping sub-slices of n size (moving windows).
// If n < 0 or len(s) == 0, the returned iterator yields nothing.
// In case len(s) < n, the returned iterator may return value only once, and length of the value may be n.
func Window[S ~[]E, E any](s S, n int) iter.Seq[S] {
	return func(yield func(S) bool) {
		if n <= 0 || len(s) == 0 {
			return
		}
		var (
			start = 0
			end   = min(n, len(s))
		)
		for {
			if !yield(s[start:end]) {
				return
			}
			start++
			end++
			if end > len(s) {
				return
			}
		}
	}
}
