package hiter

import (
	"iter"
)

func Replace[V comparable](old, new V, n int, seq iter.Seq[V]) iter.Seq[V] {
	return ReplaceFunc(old, new, n, func(i, j V) bool { return i == j }, seq)
}

func Replace2[K, V comparable](
	oldK, newK K, oldV, newV V,
	n int,
	f func(k1 K, v1 V, k2 K, v2 V) bool,
	seq iter.Seq2[K, V],
) iter.Seq2[K, V] {
	return ReplaceFunc2(
		oldK, newK, oldV, newV,
		n,
		func(k1 K, v1 V, k2 K, v2 V) bool { return k1 == k2 && v1 == v2 },
		seq,
	)
}

func ReplaceFunc[V any](old, new V, n int, f func(i, j V) bool, seq iter.Seq[V]) iter.Seq[V] {
	return func(yield func(V) bool) {
		n := n
		for v := range seq {
			if f(v, old) && n != 0 {
				if n > 0 {
					n--
				}
				v = new
			}
			if !yield(v) {
				return
			}
		}
	}
}

func ReplaceFunc2[K, V any](oldK, newK K, oldV, newV V, n int, f func(k1 K, v1 V, k2 K, v2 V) bool, seq iter.Seq2[K, V]) iter.Seq2[K, V] {
	return func(yield func(K, V) bool) {
		n := n
		for k, v := range seq {
			if f(oldK, oldV, k, v) && n != 0 {
				if n > 0 {
					n--
				}
				k, v = newK, newV
			}
			if !yield(k, v) {
				return
			}
		}
	}
}
