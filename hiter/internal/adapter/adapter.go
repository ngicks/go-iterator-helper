package adapter

import "iter"

func Concat[V any](seqs ...iter.Seq[V]) iter.Seq[V] {
	return func(yield func(V) bool) {
		for _, seq := range seqs {
			for e := range seq {
				if !yield(e) {
					return
				}
			}
		}
	}
}

func Concat2[K, V any](seqs ...iter.Seq2[K, V]) iter.Seq2[K, V] {
	return func(yield func(K, V) bool) {
		for _, seq := range seqs {
			for k, v := range seq {
				if !yield(k, v) {
					return
				}
			}
		}
	}
}

func Map[V1, V2 any](fn func(V1) V2, seq iter.Seq[V1]) iter.Seq[V2] {
	return func(yield func(V2) bool) {
		for v1 := range seq {
			if !yield(fn(v1)) {
				return
			}
		}
	}
}

func Map2[K1, K2, V1, V2 any](fn func(K1, V1) (K2, V2), seq iter.Seq2[K1, V1]) iter.Seq2[K2, V2] {
	return func(yield func(K2, V2) bool) {
		for k1, v1 := range seq {
			if !yield(fn(k1, v1)) {
				return
			}
		}
	}
}

func MergeFunc[V any](x, y iter.Seq[V], f func(V, V) int) iter.Seq[V] {
	return func(yield func(V) bool) {
		next, stop := iter.Pull(y)
		defer stop()
		v2, ok2 := next()
		for v1 := range x {
			for ok2 && f(v1, v2) > 0 {
				if !yield(v2) {
					return
				}
				v2, ok2 = next()
			}
			if !yield(v1) {
				return
			}
		}
		for ok2 {
			if !yield(v2) {
				return
			}
			v2, ok2 = next()
		}
	}
}

func Filter2[K, V any](f func(K, V) bool, seq iter.Seq2[K, V]) iter.Seq2[K, V] {
	return func(yield func(K, V) bool) {
		for k, v := range seq {
			if f(k, v) && !yield(k, v) {
				return
			}
		}
	}
}

// reduce is redefined to avoid xiter dependency.
func Reduce[Sum, V any](reducer func(Sum, V) Sum, initial Sum, seq iter.Seq[V]) Sum {
	for v := range seq {
		initial = reducer(initial, v)
	}
	return initial
}
