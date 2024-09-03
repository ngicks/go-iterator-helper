package iterable

import (
	"iter"
	"sync"

	"github.com/ngicks/go-iterator-helper/hiter"
)

var _ hiter.Iterable2[string, any] = SyncMap[string, any]{}

// SyncMap adds Iter2 method that merely calls [sync.Map.Range].
type SyncMap[K comparable, V any] struct {
	*sync.Map
}

func (s SyncMap[K, V]) Iter2() iter.Seq2[K, V] {
	return hiter.SyncMap[K, V](s.Map)
}
