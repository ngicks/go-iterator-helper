package iterable

import (
	"cmp"
	"iter"
	"maps"
	"slices"

	"github.com/ngicks/go-iterator-helper/hiter"
)

var (
	_ hiter.Iterable2[string, any] = MapAll[string, any]{}
	_ hiter.Iterable2[string, any] = MapSorted[string, any]{}
	_ hiter.Iterable2[string, any] = MapSortedFunc[map[string]any, string, any]{}
)

// MapAll adds Iter2 method to map[K]V
// that merely calling [maps.All].
type MapAll[K comparable, V any] map[K]V

func (m MapAll[K, V]) Iter2() iter.Seq2[K, V] {
	return maps.All(m)
}

// MapSorted adds Iter2 to map[K]V where K is basic ordered type.
// Iter2 takes snapshot of keys and sort it in ascending order,
// then returns an iterator over pairs of the keys and values that correspond to each key.
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

func (m MapSorted[K, V]) Reverse() MapSortedFunc[map[K]V, K, V] {
	return MapSortedFunc[map[K]V, K, V]{
		M:   m,
		Cmp: cmp.Compare[K],
		rev: true,
	}
}

// MapSortedFunc adds Iter2 to map[K]V.
// Iter2 takes snapshot of keys and sort it using the comparison function,
// then returns an iterator over pairs of the keys and values that correspond to each key.
type MapSortedFunc[M ~map[K]V, K comparable, V any] struct {
	M   M
	Cmp func(K, K) int
	rev bool
}

func (m MapSortedFunc[M, K, V]) Iter2() iter.Seq2[K, V] {
	return func(yield func(K, V) bool) {
		var keys []K
		if !m.rev {
			keys = slices.SortedFunc(maps.Keys(m.M), m.Cmp)
		} else {
			keys = slices.SortedFunc(maps.Keys(m.M), func(k1, k2 K) int { return -m.Cmp(k1, k2) })
		}
		for _, k := range keys {
			e, ok := m.M[k]
			if ok && !yield(k, e) {
				return
			}
		}
	}
}

func (m MapSortedFunc[M, K, V]) Reverse() MapSortedFunc[M, K, V] {
	return MapSortedFunc[M, K, V]{
		M:   m.M,
		Cmp: m.Cmp,
		rev: !m.rev,
	}
}
