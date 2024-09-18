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
// sentAll is false if ctx is cancelled before all values from seq are sent.
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
