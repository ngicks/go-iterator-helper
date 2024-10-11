package hiter

import "iter"

func Partition[K comparable, V any](f func(V) K, seq iter.Seq[V]) iter.Seq2[int, V] {
	return func(yield func(int, V) bool) {
		m := map[K]int{}
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
