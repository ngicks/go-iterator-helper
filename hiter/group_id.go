package hiter

import "iter"

// WithGroupId returns an iterator over pair of unique index for the group id and value from seq.
// f determines unique group identifiers that corresponds the input value.
// WithGroupId buffers all group ids f returns, then converts its appearance order.
func WithGroupId[V any, I comparable](f func(V) I, seq iter.Seq[V]) iter.Seq2[int, V] {
	return func(yield func(int, V) bool) {
		m := map[I]int{}
		for v := range seq {
			k := f(v)
			order, ok := m[k]
			if !ok {
				order = len(m)
				m[k] = order
			}
			if !yield(order, v) {
				return
			}
		}
	}
}
