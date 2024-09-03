package collection

import (
	"iter"

	"github.com/ngicks/go-iterator-helper/hiter"
)

// reduce is redefined to avoid xiter dependency.
func reduce[Sum, V any](seq iter.Seq[V], reducer func(Sum, V) Sum, initial Sum) Sum {
	for v := range seq {
		initial = reducer(initial, v)
	}
	return initial
}

func SumOf[T any, E hiter.Numeric](seq iter.Seq[T], selector func(ele T) E) E {
	return reduce(
		seq,
		func(e E, t T) E { return e + selector(t) },
		*new(E),
	)
}
