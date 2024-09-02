package iterable

import (
	"iter"
	"sync"

	"github.com/ngicks/go-iterator-helper/ih"
)

var _ ih.Iterable2[string, any] = SyncMap[string, any]{}

// SyncMap adds Iter2 method that merely calls [sync.Map.Range].
type SyncMap[K comparable, V any] struct {
	*sync.Map
}

func (s SyncMap[K, V]) Iter2() iter.Seq2[K, V] {
	return ih.SyncMap[K, V](s.Map)
}
