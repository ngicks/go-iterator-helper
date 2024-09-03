package hiter

import "iter"

// Window returns an iterator over overlapping sub-slices of n size (moving windows).
func Window[S ~[]E, E any](s S, n int) iter.Seq[S] {
	return func(yield func(S) bool) {
		if n <= 0 {
			return
		}
		var (
			start = 0
			end   = n
			done  = false
		)
		if end > len(s) {
			end = len(s)
			done = true
		}
		for {
			if !yield(s[start:end]) {
				return
			}
			if done {
				return
			}
			start++
			end++
			if end > len(s) {
				end = len(s)
			}
			if end == len(s) {
				done = true
			}
		}
	}
}
