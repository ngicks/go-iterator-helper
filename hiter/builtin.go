package hiter

import (
	"iter"
)

// MapKeys returns an iterator over pairs of keys and values retrieved from a map using the keys.
// The keys iterator determines the order of access.
func MapKeys[M ~map[K]V, K comparable, V any](m M, keys iter.Seq[K]) iter.Seq2[K, V] {
	return func(yield func(K, V) bool) {
		for k := range keys {
			if !yield(k, m[k]) {
				return
			}
		}
	}
}
