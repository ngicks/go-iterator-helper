package hiter

import (
	"iter"
)

func ReduceGroup[K comparable, V, Sum any](seq iter.Seq2[K, V], reducer func(accumulator Sum, current V) Sum, initial Sum) map[K]Sum {
	m := make(map[K]Sum)
	for k, v := range seq {
		if _, ok := m[k]; !ok {
			m[k] = initial
		}
		m[k] = reducer(m[k], v)
	}
	return m
}
