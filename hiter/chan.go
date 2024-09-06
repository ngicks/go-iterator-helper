package hiter

import (
	"context"
	"iter"
)

// Chan returns an iterator over ch.
// Either cancelling ctx or closing ch stops iteration.
// ctx is allowed to be nil.
func Chan[V any](ctx context.Context, ch <-chan V) iter.Seq[V] {
	return func(yield func(V) bool) {
		if ctx == nil {
			ctx = context.Background()
		}
		for {
			select {
			case <-ctx.Done():
				return
			case v, ok := <-ch:
				if !ok || !yield(v) {
					return
				}
			}
		}
	}
}
