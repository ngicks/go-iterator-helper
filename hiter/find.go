package hiter

import (
	"iter"
)

// Find iterates over seq and the first value equals to v with its count.
// The count is 0-indexed and -1 if seq does not have value equals v.
func Find[V comparable](v V, seq iter.Seq[V]) (V, int) {
	var i int
	for vv := range seq {
		if vv == v {
			return vv, i
		}
		i++
	}
	return *new(V), -1
}

// FindFunc iterates over seq and the first value satisfying fn(v) with its count.
// The count is 0-indexed and -1 if seq does not have value that satisfies fn.
func FindFunc[V any](f func(V) bool, seq iter.Seq[V]) (V, int) {
	var i int
	for vv := range seq {
		if f(vv) {
			return vv, i
		}
		i++
	}
	return *new(V), -1
}

// Find2 iterates over seq and the first k-v pair equals to k and v with its count.
// The count is 0-indexed and -1 if seq does not have k-v pairs equals the input.
func Find2[K, V comparable](k K, v V, seq iter.Seq2[K, V]) (K, V, int) {
	var i int
	for kk, vv := range seq {
		if kk == k && vv == v {
			return kk, vv, i
		}
		i++
	}
	return *new(K), *new(V), -1
}

// FindFunc2 iterates over seq and the first k-v pair satisfying fn(v) with its count.
// The count is 0-indexed and -1 if seq does not have k-v pairs that satisfies fn.
func FindFunc2[K, V any](fn func(K, V) bool, seq iter.Seq2[K, V]) (K, V, int) {
	var i int
	for kk, vv := range seq {
		if fn(kk, vv) {
			return kk, vv, i
		}
		i++
	}
	return *new(K), *new(V), -1
}

// FindLast is like [Find] but returns the final occurrence of v.
// Unlike [Find], FindLast always fully consumes seq.
func FindLast[V comparable](v V, seq iter.Seq[V]) (found V, idx int) {
	idx = -1
	var i int
	for vv := range seq {
		if vv == v {
			found = vv
			idx = i
		}
		i++
	}
	return found, idx
}

// FindLastFunc is like [FindFunc] but returns the final occurrence of v.
func FindLastFunc[V any](fn func(V) bool, seq iter.Seq[V]) (found V, idx int) {
	idx = -1
	var i int
	for vv := range seq {
		if fn(vv) {
			found = vv
			idx = i
		}
		i++
	}
	return found, idx
}

// FindLast2 is like [Find2] but returns the final occurrence of k-v pair.
func FindLast2[K, V comparable](k K, v V, seq iter.Seq2[K, V]) (foundK K, foundV V, idx int) {
	idx = -1
	var i int
	for kk, vv := range seq {
		if kk == k && vv == v {
			foundK = kk
			foundV = vv
			idx = i
		}
		i++
	}
	return foundK, foundV, idx
}

// FindLastFunc2 is like [FindFunc2] but returns the final occurrence of k-v pair.
func FindLastFunc2[K, V any](fn func(K, V) bool, seq iter.Seq2[K, V]) (foundK K, foundV V, idx int) {
	idx = -1
	var i int
	for kk, vv := range seq {
		if fn(kk, vv) {
			foundK = kk
			foundV = vv
			idx = i
		}
		i++
	}
	return foundK, foundV, idx
}

// Contains reports whether v is present in seq.
func Contains[V comparable](v V, seq iter.Seq[V]) bool {
	_, idx := Find(v, seq)
	return idx >= 0
}

// ContainsFunc reports whether at least one element v of s satisfies f(v).
func ContainsFunc[V any](f func(V) bool, seq iter.Seq[V]) bool {
	_, idx := FindFunc(f, seq)
	return idx >= 0
}

// Contains2 reports whether k-v pair is present in seq.
func Contains2[K, V comparable](k K, v V, seq iter.Seq2[K, V]) bool {
	_, _, idx := Find2(k, v, seq)
	return idx >= 0
}

// ContainsFunc2 reports whether at least one k-v pair of s satisfies f(k, v).
func ContainsFunc2[K, V any](f func(K, V) bool, seq iter.Seq2[K, V]) bool {
	_, _, idx := FindFunc2(f, seq)
	return idx >= 0
}
