package ih

import "iter"

var _ Iterable2[any, any] = KeyValues[any, any](nil)

// KeyValues adds the Iter2 method to slice of KeyValue-s.
type KeyValues[K, V any] []KeyValue[K, V]

func (v KeyValues[K, V]) Iter2() iter.Seq2[K, V] {
	return func(yield func(K, V) bool) {
		for _, kv := range v {
			if !yield(kv.K, kv.V) {
				return
			}
		}
	}
}
