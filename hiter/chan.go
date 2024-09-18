package hiter

import (
	"context"
	"iter"
)

// Chan returns an iterator over ch.
// Either cancelling ctx or closing ch stops iteration.
func Chan[V any](ctx context.Context, ch <-chan V) iter.Seq[V] {
	return func(yield func(V) bool) {
		for {
			select {
			case <-ctx.Done():
				return
			default:
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
}

// ChanSend sends values from seq to c.
// It unblocks after either cancelling ctx or consuming all values from seq.
// sentAll is true only when all values are sent to c.
func ChanSend[V any](ctx context.Context, c chan<- V, seq iter.Seq[V]) (v V, sentAll bool) {
	for v := range seq {
		select {
		case <-ctx.Done():
			return v, false
		case c <- v:
		}
	}
	return *new(V), true
}
