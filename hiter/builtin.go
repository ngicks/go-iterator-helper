package hiter

import (
	"cmp"
	"iter"
	"maps"
	"slices"
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

// MapSorted returns an iterator over key-value pairs from m.
// The iterator takes snapshot of keys from m at invocation and sorts keys in ascending order,
// which determines the order of paris.
//
// Larger the map is, more the memory is consumed.
func MapSorted[M ~map[K]V, K cmp.Ordered, V any](m M) iter.Seq2[K, V] {
	return MapSortedFunc(m, cmp.Compare)
}

// MapSorted is like [MapSorted] but uses comparison function f to sort keys.
func MapSortedFunc[M ~map[K]V, K comparable, V any](m M, f func(i, j K) int) iter.Seq2[K, V] {
	return func(yield func(K, V) bool) {
		keys := slices.SortedFunc(maps.Keys(m), f)
		for _, k := range keys {
			if !yield(k, m[k]) {
				return
			}
		}
	}
}
