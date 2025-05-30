package hiter

import (
	"cmp"
	"iter"
)

// Concat returns an iterator over the concatenation of the sequences.
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

// Concat2 returns an iterator over the concatenation of the sequences.
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

// Equal return true if 2 seqs yields same sequences.
func Equal[V comparable](x, y iter.Seq[V]) bool {
	for z := range Zip(x, y) {
		if z.L.Ok != z.R.Ok || z.L.V != z.R.V {
			return false
		}
	}
	return true
}

// Equal2 return true if 2 seqs yields same sequences.
func Equal2[K, V comparable](x, y iter.Seq2[K, V]) bool {
	for z := range Zip2(x, y) {
		if z.L.Ok != z.R.Ok || z.L.V.K != z.R.V.K || z.L.V.V != z.R.V.V {
			return false
		}
	}
	return true
}

// EqualFunc reports whether two sequences are equal using an equality function on each pair of elements.
// If the lengths are different, EqualFunc returns false.
// Otherwise, the elements are compared in increasing index order,
// and the comparison stops at the first index for which eq returns false.
func EqualFunc[V1, V2 any](x iter.Seq[V1], y iter.Seq[V2], f func(V1, V2) bool) bool {
	for z := range Zip(x, y) {
		if z.L.Ok != z.R.Ok || !f(z.L.V, z.R.V) {
			return false
		}
	}
	return true
}

// EqualFunc2 reports whether two sequences are equal using an equality function on each pair of elements.
// If the lengths are different, EqualFunc returns false.
// Otherwise, the elements are compared in increasing index order,
// and the comparison stops at the first index for which eq returns false.
func EqualFunc2[K1, V1, K2, V2 any](x iter.Seq2[K1, V1], y iter.Seq2[K2, V2], f func(K1, V1, K2, V2) bool) bool {
	for z := range Zip2(x, y) {
		if z.L.Ok != z.R.Ok || !f(z.L.V.K, z.L.V.V, z.R.V.K, z.R.V.V) {
			return false
		}
	}
	return true
}

// Filter returns an iterator over seq that only includes
// the values v for which f(v) is true.
func Filter[V any](f func(V) bool, seq iter.Seq[V]) iter.Seq[V] {
	return func(yield func(V) bool) {
		for v := range seq {
			if f(v) && !yield(v) {
				return
			}
		}
	}
}

// Filter2 returns an iterator over seq that only includes
// the pairs k, v for which f(k, v) is true.
func Filter2[K, V any](f func(K, V) bool, seq iter.Seq2[K, V]) iter.Seq2[K, V] {
	return func(yield func(K, V) bool) {
		for k, v := range seq {
			if f(k, v) && !yield(k, v) {
				return
			}
		}
	}
}

// Limit returns an iterator over seq that stops after n values.
func Limit[V any](n int, seq iter.Seq[V]) iter.Seq[V] {
	return func(yield func(V) bool) {
		n := n // no state in the iterator
		if n <= 0 {
			return
		}
		for v := range seq {
			if !yield(v) {
				return
			}
			if n--; n <= 0 {
				break
			}
		}
	}
}

// Limit2 returns an iterator over seq that stops after n key-value pairs.
func Limit2[K, V any](n int, seq iter.Seq2[K, V]) iter.Seq2[K, V] {
	return func(yield func(K, V) bool) {
		n := n
		if n <= 0 {
			return
		}
		for k, v := range seq {
			if !yield(k, v) {
				return
			}
			if n--; n <= 0 {
				break
			}
		}
	}
}

// Map returns an iterator over f applied to seq.
func Map[V1, V2 any](f func(V1) V2, seq iter.Seq[V1]) iter.Seq[V2] {
	return func(yield func(V2) bool) {
		for in := range seq {
			if !yield(f(in)) {
				return
			}
		}
	}
}

// Map2 returns an iterator over f applied to seq.
func Map2[K1, V1, K2, V2 any](f func(K1, V1) (K2, V2), seq iter.Seq2[K1, V1]) iter.Seq2[K2, V2] {
	return func(yield func(K2, V2) bool) {
		for k, v := range seq {
			if !yield(f(k, v)) {
				return
			}
		}
	}
}

// Merge merges two sequences of ordered values.
// Values appear in the output once for each time they appear in x
// and once for each time they appear in y.
// If the two input sequences are not ordered,
// the output sequence will not be ordered,
// but it will still contain every value from x and y exactly once.
//
// Merge is equivalent to calling MergeFunc with cmp.Compare[V]
// as the ordering function.
func Merge[V cmp.Ordered](x, y iter.Seq[V]) iter.Seq[V] {
	return MergeFunc(x, y, cmp.Compare[V])
}

// MergeFunc merges two sequences of values ordered by the function f.
// Values appear in the output once for each time they appear in x
// and once for each time they appear in y.
// When equal values appear in both sequences,
// the output contains the values from x before the values from y.
// If the two input sequences are not ordered by f,
// the output sequence will not be ordered by f,
// but it will still contain every value from x and y exactly once.
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

// Merge2 merges two sequences of key-value pairs ordered by their keys.
// Pairs appear in the output once for each time they appear in x
// and once for each time they appear in y.
// If the two input sequences are not ordered by their keys,
// the output sequence will not be ordered by its keys,
// but it will still contain every pair from x and y exactly once.
//
// Merge2 is equivalent to calling MergeFunc2 with cmp.Compare[K]
// as the ordering function.
func Merge2[K cmp.Ordered, V any](x, y iter.Seq2[K, V]) iter.Seq2[K, V] {
	return MergeFunc2(x, y, cmp.Compare[K])
}

// MergeFunc2 merges two sequences of key-value pairs ordered by the function f.
// Pairs appear in the output once for each time they appear in x
// and once for each time they appear in y.
// When pairs with equal keys appear in both sequences,
// the output contains the pairs from x before the pairs from y.
// If the two input sequences are not ordered by f,
// the output sequence will not be ordered by f,
// but it will still contain every pair from x and y exactly once.
func MergeFunc2[K, V any](x, y iter.Seq2[K, V], f func(K, K) int) iter.Seq2[K, V] {
	return func(yield func(K, V) bool) {
		next, stop := iter.Pull2(y)
		defer stop()
		k2, v2, ok2 := next()
		for k1, v1 := range x {
			for ok2 && f(k1, k2) > 0 {
				if !yield(k2, v2) {
					return
				}
				k2, v2, ok2 = next()
			}
			if !yield(k1, v1) {
				return
			}
		}
		for ok2 {
			if !yield(k2, v2) {
				return
			}
			k2, v2, ok2 = next()
		}
	}
}

// Reduce combines the values in seq using f.
// For each value v in seq, it updates sum = f(sum, v)
// and then returns the final sum.
// For example, if iterating over seq yields v1, v2, v3,
// Reduce returns f(f(f(sum, v1), v2), v3).
func Reduce[Sum, V any](f func(Sum, V) Sum, sum Sum, seq iter.Seq[V]) Sum {
	for v := range seq {
		sum = f(sum, v)
	}
	return sum
}

// Reduce2 combines the values in seq using f.
// For each pair k, v in seq, it updates sum = f(sum, k, v)
// and then returns the final sum.
// For example, if iterating over seq yields (k1, v1), (k2, v2), (k3, v3)
// Reduce returns f(f(f(sum, k1, v1), k2, v2), k3, v3).
func Reduce2[Sum, K, V any](f func(Sum, K, V) Sum, sum Sum, seq iter.Seq2[K, V]) Sum {
	for k, v := range seq {
		sum = f(sum, k, v)
	}
	return sum
}

type Option[V any] struct {
	V  V
	Ok bool
}

// A Zipped is a pair of zipped values, one of which may be missing,
// drawn from two different sequences.
type Zipped[V1, V2 any] struct {
	L Option[V1]
	R Option[V2]
}

func Zip[V1, V2 any](x iter.Seq[V1], y iter.Seq[V2]) iter.Seq[Zipped[V1, V2]] {
	return func(yield func(z Zipped[V1, V2]) bool) {
		next, stop := iter.Pull(y)
		defer stop()
		v2, ok2 := next()
		for v1 := range x {
			if !yield(Zipped[V1, V2]{Option[V1]{v1, true}, Option[V2]{v2, ok2}}) {
				return
			}
			v2, ok2 = next()
		}
		var zv1 V1
		for ok2 {
			if !yield(Zipped[V1, V2]{Option[V1]{zv1, false}, Option[V2]{v2, ok2}}) {
				return
			}
			v2, ok2 = next()
		}
	}
}

func Zip2[K1, V1, K2, V2 any](x iter.Seq2[K1, V1], y iter.Seq2[K2, V2]) iter.Seq[Zipped[KeyValue[K1, V1], KeyValue[K2, V2]]] {
	return func(yield func(z Zipped[KeyValue[K1, V1], KeyValue[K2, V2]]) bool) {
		next, stop := iter.Pull2(y)
		defer stop()
		k2, v2, ok2 := next()
		for k1, v1 := range x {
			if !yield(Zipped[KeyValue[K1, V1], KeyValue[K2, V2]]{
				Option[KeyValue[K1, V1]]{KeyValue[K1, V1]{k1, v1}, true},
				Option[KeyValue[K2, V2]]{KeyValue[K2, V2]{k2, v2}, ok2},
			}) {
				return
			}
			k2, v2, ok2 = next()
		}
		var zk1 K1
		var zv1 V1
		for ok2 {
			if !yield(Zipped[KeyValue[K1, V1], KeyValue[K2, V2]]{
				Option[KeyValue[K1, V1]]{KeyValue[K1, V1]{zk1, zv1}, false},
				Option[KeyValue[K2, V2]]{KeyValue[K2, V2]{k2, v2}, ok2},
			}) {
				return
			}
			k2, v2, ok2 = next()
		}
	}
}
