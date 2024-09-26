package hiter

import "iter"

// Empty returns an iterator over nothing.
func Empty[V any]() iter.Seq[V] {
	return func(yield func(V) bool) {}
}

// Empty2 returns an iterator over nothing.
func Empty2[K, V any]() iter.Seq2[K, V] {
	return func(yield func(K, V) bool) {}
}
