package hiter

import (
	"iter"
)

// reduce is redefined to avoid xiter dependency.
func reduce[Sum, V any](seq iter.Seq[V], reducer func(Sum, V) Sum, initial Sum) Sum {
	for v := range seq {
		initial = reducer(initial, v)
	}
	return initial
}

// as per https://go.dev/ref/spec#arithmetic_operators
type Summable interface {
	~int | ~int8 | ~int16 | ~int32 | ~int64 |
		~uint | ~uint8 | ~uint16 | ~uint32 | ~uint64 | ~uintptr |
		~float32 | ~float64 |
		~complex64 | ~complex128 |
		~string
}

func Sum[S Summable](seq iter.Seq[S]) S {
	return reduce(
		seq,
		func(e S, t S) S { return e + t },
		*new(S),
	)
}

func SumOf[V any, S Summable](selector func(ele V) S, seq iter.Seq[V]) S {
	return reduce(
		seq,
		func(e S, t V) S { return e + selector(t) },
		*new(S),
	)
}
