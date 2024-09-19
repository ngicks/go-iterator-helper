package hiter

import "iter"

func RunningReduce[V, Sum any](reducer func(accumulator Sum, current V, i int) Sum, initial Sum, seq iter.Seq[V]) iter.Seq[Sum] {
	return func(yield func(Sum) bool) {
		var i int
		for v := range seq {
			initial = reducer(initial, v, i)
			i++
			if !yield(initial) {
				return
			}
		}
	}
}
