package hiter

import (
	"iter"
)

// Replace is [ReplaceFunc] where matcher f tests if k equals old.
func Replace[V comparable](old, new V, n int, seq iter.Seq[V]) iter.Seq[V] {
	return ReplaceFunc(
		func(v V) bool { return v == old },
		new,
		n,
		seq,
	)
}

// Replace2 is [ReplaceFunc2] where matcher f tests if k-v pair equals pair of oldK, oldV.
func Replace2[K, V comparable](
	oldK K, oldV V, newK K, newV V,
	n int,
	seq iter.Seq2[K, V],
) iter.Seq2[K, V] {
	return ReplaceFunc2(
		func(k K, v V) bool { return k == oldK && v == oldV },
		newK, newV,
		n,
		seq,
	)
}

// ReplaceFunc returns an iterator over seq but the first n values that matches f is replaced with new.
// If n < 0, there is no limit on the number of replacements.
func ReplaceFunc[V any](f func(V) bool, new V, n int, seq iter.Seq[V]) iter.Seq[V] {
	if n == 0 {
		// no allocation
		return seq
	}
	return func(yield func(V) bool) {
		n := n
		for v := range seq {
			if f(v) && n != 0 {
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

// ReplaceFunc2 returns an iterator over seq but the first n key-value pairs that matches f is replaced with new.
// If n < 0, there is no limit on the number of replacements.
func ReplaceFunc2[K, V any](
	f func(K, V) bool,
	newK K, newV V,
	n int,
	seq iter.Seq2[K, V],
) iter.Seq2[K, V] {
	if n == 0 {
		// no allocation
		return seq
	}
	return func(yield func(K, V) bool) {
		n := n
		for k, v := range seq {
			if f(k, v) && n != 0 {
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
