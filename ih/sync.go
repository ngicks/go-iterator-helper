package ih

import (
	"iter"
	"sync"
)

// SyncMap returns an iterator over m.
// Breaking Seq2 may stop producing more data, however it may still be O(N).
func SyncMap[K, V any](m *sync.Map) iter.Seq2[K, V] {
	return func(yield func(K, V) bool) {
		m.Range(func(key, value any) bool {
			return yield(key.(K), value.(V))
		})
	}
}
