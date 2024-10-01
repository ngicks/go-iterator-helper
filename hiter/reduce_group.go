package hiter

import (
	"iter"
)

// ReduceGroup sums up values from seq for every unique keys separately,
// then collects pairs of key and converted value into a new map and returns it.
// initial is passed once for every unique keys to reducer.
func ReduceGroup[K comparable, V, Sum any](
	reducer func(accumulator Sum, current V) Sum,
	initial Sum,
	seq iter.Seq2[K, V],
) map[K]Sum {
	return InsertReduceGroup(make(map[K]Sum), reducer, initial, seq)
}

// InsertReduceGroup is like [ReduceGroup] but inserts result to m.
func InsertReduceGroup[Map ~map[K]Sum, K comparable, V, Sum any](
	m Map,
	reducer func(accumulator Sum, current V) Sum,
	initial Sum,
	seq iter.Seq2[K, V],
) map[K]Sum {
	for k, v := range seq {
		if _, ok := m[k]; !ok {
			m[k] = initial
		}
		m[k] = reducer(m[k], v)
	}
	return m
}
