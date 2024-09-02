package ih

import "iter"

// Chan returns an iterator over ch.
// Only closing ch stops the iterator.
func Chan[V any](ch <-chan V, f func()) iter.Seq[V] {
	return func(yield func(V) bool) {
		defer f()
		for v := range ch {
			if !yield(v) {
				return
			}
		}
	}
}
