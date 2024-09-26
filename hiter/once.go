package hiter

import "iter"

// Once adapts a single value as an iterator;
// the iterator yields v and stops.
func Once[V any](v V) iter.Seq[V] {
	return func(yield func(V) bool) {
		yield(v)
	}
}

// Once2 adapts a single k-v pair as an iterator;
// the iterator yields k, v and stops.
func Once2[K, V any](k K, v V) iter.Seq2[K, V] {
	return func(yield func(K, V) bool) {
		yield(k, v)
	}
}
