package ih

import "iter"

// Chunk returns an iterator over non overlapping sub-slices of n size.
func Chunk[S ~[]E, E any](s S, n int) iter.Seq[S] {
	return func(yield func(S) bool) {
		if n <= 0 {
			return
		}
		var cut S
		for {
			if len(s) >= n {
				cut, s = s[:n], s[n:]
			} else {
				cut, s = s, nil
			}
			if len(cut) == 0 {
				return
			}
			if !yield(cut) {
				return
			}
		}
	}
}
