package hiter

import (
	"iter"

	"github.com/ngicks/go-iterator-helper/hiter/internal/adapter"
)

// Assert returns an iterator over seq but each value is type-asserted to be type V.
func Assert[V any](seq iter.Seq[any]) iter.Seq[V] {
	return adapter.Map(func(v any) V { return v.(V) }, seq)
}

// Assert2 returns an iterator over seq but each key-value pair is type-asserted to be type K and V respectively.
func Assert2[K, V any](seq iter.Seq2[any, any]) iter.Seq2[K, V] {
	return adapter.Map2(func(k, v any) (K, V) { return k.(K), v.(V) }, seq)
}
