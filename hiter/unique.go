package hiter

import (
	"iter"
	"unique"
)

func Unique[V comparable](seq iter.Seq[V]) iter.Seq[V] {
	return func(yield func(V) bool) {
		handles := map[unique.Handle[V]]bool{}
		for v := range seq {
			handle := unique.Make(v)
			if handles[handle] {
				continue
			}
			handles[handle] = true
			if !yield(v) {
				return
			}
		}
	}
}

func Unique2[K, V comparable](seq iter.Seq2[K, V]) iter.Seq2[K, V] {
	return func(yield func(K, V) bool) {
		handles := map[KeyValue[unique.Handle[K], unique.Handle[V]]]bool{}
		for k, v := range seq {
			handle := KeyValue[unique.Handle[K], unique.Handle[V]]{unique.Make(k), unique.Make(v)}
			if handles[handle] {
				continue
			}
			handles[handle] = true
			if !yield(k, v) {
				return
			}
		}
	}
}
