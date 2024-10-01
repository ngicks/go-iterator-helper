package hiter

import (
	"cmp"
	"iter"
)

func Min[V cmp.Ordered](seq iter.Seq[V]) V {
	return MinFunc(cmp.Compare, seq)
}

func MinFunc[V any](fn func(i, j V) int, seq iter.Seq[V]) V {
	var (
		min   V
		first = true
	)
	for v := range seq {
		if first {
			min = v
			first = false
		} else {
			if fn(v, min) < 0 {
				min = v
			}
		}
	}
	return min
}

func Max[V cmp.Ordered](seq iter.Seq[V]) V {
	return MaxFunc(cmp.Compare, seq)
}

func MaxFunc[V any](fn func(i, j V) int, seq iter.Seq[V]) V {
	var (
		max   V
		first = true
	)
	for v := range seq {
		if first {
			max = v
			first = false
		} else {
			if fn(v, max) > 0 {
				max = v
			}
		}
	}
	return max
}
