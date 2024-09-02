package iterable

import (
	"cmp"
	"iter"
	"maps"
	"slices"

	"github.com/ngicks/go-iterator-helper/ih"
)

var (
	_ ih.Iterable2[string, any] = MapAll[string, any]{}
	_ ih.Iterable2[string, any] = MapSorted[string, any]{}
	_ ih.Iterable2[string, any] = MapSortedFunc[map[string]any, string, any]{}
)

// MapAll adds Iter2 method to map[K]V
// that merely calling maps.All.
type MapAll[K comparable, V any] map[K]V

func (m MapAll[K, V]) Iter2() iter.Seq2[K, V] {
	return maps.All(m)
}

// MapSorted adds Iter2 to map[K]V where K is basic ordered type.
// Iter2 takes snapshot of keys and sort it in ascending order,
// then returns an iterator over pairs of the keys and values that corresponds to the key.
type MapSorted[K cmp.Ordered, V any] map[K]V

func (m MapSorted[K, V]) Iter2() iter.Seq2[K, V] {
	return func(yield func(K, V) bool) {
		for _, k := range slices.Sorted(maps.Keys(m)) {
			e, ok := m[k]
			if ok && !yield(k, e) {
				return
			}
		}
	}
}

// MapSortedFunc adds Iter2 to map[K]V.
// Iter2 takes snapshot of keys and sort it using the comparison function,
// then returns an iterator over pairs of the keys and values that corresponds to the key.
type MapSortedFunc[M ~map[K]V, K comparable, V any] struct {
	M   M
	Cmp func(K, K) int
}

func (m MapSortedFunc[M, K, V]) Iter2() iter.Seq2[K, V] {
	return func(yield func(K, V) bool) {
		for _, k := range slices.SortedFunc(maps.Keys(m.M), m.Cmp) {
			e, ok := m.M[k]
			if ok && !yield(k, e) {
				return
			}
		}
	}
}
