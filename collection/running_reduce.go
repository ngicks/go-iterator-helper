package collection

import "iter"

func RunningReduce[V, Sum any](seq iter.Seq[V], reducer func(accumulator Sum, current V, i int) Sum, initial Sum) iter.Seq[Sum] {
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
