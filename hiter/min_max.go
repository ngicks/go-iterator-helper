package hiter

import (
	"cmp"
	"iter"
)

// Min returns the minimum value of seq.
// Min returns the zero value if seq is empty.
func Min[V cmp.Ordered](seq iter.Seq[V]) V {
	return MinFunc(cmp.Compare, seq)
}

// MinFunc returns the minimum value of seq using comparison function fn.
// fn must return -1 if x is less than y, 0 if x equals y, and +1 x is greater than y (as [cmp.Compare] does.)
// MinFunc returns the zero value if seq is empty.
func MinFunc[V any](fn func(x, y V) int, seq iter.Seq[V]) V {
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

// Max returns the maximum value of seq.
// Max returns the zero value if seq is empty.
func Max[V cmp.Ordered](seq iter.Seq[V]) V {
	return MaxFunc(cmp.Compare, seq)
}

// MaxFunc returns the maximum value of seq using comparison function fn.
// fn must return -1 if x is less than y, 0 if x equals y, and +1 x is greater than y (as [cmp.Compare] does.)
// MaxFunc returns the zero value if seq is empty.
func MaxFunc[V any](fn func(x, y V) int, seq iter.Seq[V]) V {
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
