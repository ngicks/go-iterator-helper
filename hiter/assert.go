package hiter

import (
	"iter"
	"reflect"
)

func mapIter[V1, V2 any](fn func(V1) V2, seq iter.Seq[V1]) iter.Seq[V2] {
	return func(yield func(V2) bool) {
		for v1 := range seq {
			if !yield(fn(v1)) {
				return
			}
		}
	}
}

func mapIter2[K1, K2, V1, V2 any](fn func(K1, V1) (K2, V2), seq iter.Seq2[K1, V1]) iter.Seq2[K2, V2] {
	return func(yield func(K2, V2) bool) {
		for k1, v1 := range seq {
			if !yield(fn(k1, v1)) {
				return
			}
		}
	}
}

// AssertValue returns an iterator over seq but each value returned by [reflect.Value.Interface] is type-asserted to be type V.
func AssertValue[V any](seq iter.Seq[reflect.Value]) iter.Seq[V] {
	return mapIter(func(v reflect.Value) V { return v.Interface().(V) }, seq)
}

// Assert2 returns an iterator over seq but internal values returned by [reflect.Value.Interface] of each key-value pair
// are type-asserted to be type K and V respectively.
func AssertValue2[K, V any](seq iter.Seq2[reflect.Value, reflect.Value]) iter.Seq2[K, V] {
	return mapIter2(func(k, v reflect.Value) (K, V) { return k.Interface().(K), v.Interface().(V) }, seq)
}

// Assert returns an iterator over seq but each value is type-asserted to be type V.
func Assert[V any](seq iter.Seq[any]) iter.Seq[V] {
	return mapIter(func(v any) V { return v.(V) }, seq)
}

// Assert2 returns an iterator over seq but each key-value pair is type-asserted to be type K and V respectively.
func Assert2[K, V any](seq iter.Seq2[any, any]) iter.Seq2[K, V] {
	return mapIter2(func(k, v any) (K, V) { return k.(K), v.(V) }, seq)
}
