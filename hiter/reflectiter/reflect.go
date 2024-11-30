package reflectiter

import (
	"iter"
	"reflect"

	"github.com/ngicks/go-iterator-helper/hiter/internal/adapter"
)

// AssertValue returns an iterator over seq but each value returned by [reflect.Value.Interface] is type-asserted to be type V.
func AssertValue[V any](seq iter.Seq[reflect.Value]) iter.Seq[V] {
	return adapter.Map(func(v reflect.Value) V { return v.Interface().(V) }, seq)
}

// Assert2 returns an iterator over seq but internal values returned by [reflect.Value.Interface] of each key-value pair
// are type-asserted to be type K and V respectively.
func AssertValue2[K, V any](seq iter.Seq2[reflect.Value, reflect.Value]) iter.Seq2[K, V] {
	return adapter.Map2(func(k, v reflect.Value) (K, V) { return k.Interface().(K), v.Interface().(V) }, seq)
}

// Seq wraps [reflect.Value.Seq] with [AssertValue].
func Seq[V any](rv reflect.Value) iter.Seq[V] {
	return AssertValue[V](rv.Seq())
}

// Seq2 wraps [reflect.Value.Seq2] with [AssertValue2].
func Seq2[K, V any](rv reflect.Value) iter.Seq2[K, V] {
	return AssertValue2[K, V](rv.Seq2())
}
