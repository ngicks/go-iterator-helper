package hiter

import (
	"iter"
	"reflect"

	"github.com/ngicks/go-iterator-helper/hiter/internal/adapter"
)

// AssertValue returns an iterator over seq but each value returned by [reflect.Value.Interface] is type-asserted to be type V.
func AssertValue[V any](seq iter.Seq[reflect.Value]) iter.Seq[V] {
	return adapter.MapIter(func(v reflect.Value) V { return v.Interface().(V) }, seq)
}

// Assert2 returns an iterator over seq but internal values returned by [reflect.Value.Interface] of each key-value pair
// are type-asserted to be type K and V respectively.
func AssertValue2[K, V any](seq iter.Seq2[reflect.Value, reflect.Value]) iter.Seq2[K, V] {
	return adapter.MapIter2(func(k, v reflect.Value) (K, V) { return k.Interface().(K), v.Interface().(V) }, seq)
}

// Assert returns an iterator over seq but each value is type-asserted to be type V.
func Assert[V any](seq iter.Seq[any]) iter.Seq[V] {
	return adapter.MapIter(func(v any) V { return v.(V) }, seq)
}

// Assert2 returns an iterator over seq but each key-value pair is type-asserted to be type K and V respectively.
func Assert2[K, V any](seq iter.Seq2[any, any]) iter.Seq2[K, V] {
	return adapter.MapIter2(func(k, v any) (K, V) { return k.(K), v.(V) }, seq)
}
