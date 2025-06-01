package hiter

import (
	"cmp"
	"iter"
	"maps"
	"slices"
)

// MapsKeys returns an iterator over pairs of keys and values retrieved from a map using the keys.
// The keys iterator determines the order of access.
func MapsKeys[M ~map[K]V, K comparable, V any](m M, keys iter.Seq[K]) iter.Seq2[K, V] {
	return func(yield func(K, V) bool) {
		for k := range keys {
			if !yield(k, m[k]) {
				return
			}
		}
	}
}

// MapsSorted returns an iterator over key-value pairs from m.
// The iterator takes snapshot of keys from m at invocation and sorts keys in ascending order,
// which determines the order of paris.
//
// Larger the map is, more the memory is consumed.
func MapsSorted[M ~map[K]V, K cmp.Ordered, V any](m M) iter.Seq2[K, V] {
	return MapsSortedFunc(m, cmp.Compare)
}

// MapSorted is like [MapsSorted] but uses comparison function f to sort keys.
func MapsSortedFunc[M ~map[K]V, K comparable, V any](m M, f func(i, j K) int) iter.Seq2[K, V] {
	return func(yield func(K, V) bool) {
		keys := slices.SortedFunc(maps.Keys(m), f)
		for _, k := range keys {
			if !yield(k, m[k]) {
				return
			}
		}
	}
}

// MapsOverlay return an iterator over key-value pairs from maps mm in overlaying manner.
// Every unique keys appear in the output only once for each.
// If multiple maps has same key, value in right most map is chosen.
func MapsOverlay[M ~map[K]V, K comparable, V any](mm ...M) iter.Seq2[K, V] {
	return func(yield func(K, V) bool) {
		seen := map[K]struct{}{}
		for _, m := range slices.Backward(mm) {
			for k, v := range m {
				if _, ok := seen[k]; ok {
					continue
				}
				seen[k] = struct{}{}
				if !yield(k, v) {
					return
				}
			}
		}
	}
}
