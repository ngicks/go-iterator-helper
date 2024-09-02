package iterable

import (
	"iter"

	"github.com/ngicks/go-iterator-helper/ih"
)

var _ ih.IntoIterable[any] = Chan[any](nil)

// Chan adds IntoIter method to a receive only channel.
type Chan[V any] <-chan V

func (c Chan[V]) IntoIter() iter.Seq[V] {
	return func(yield func(V) bool) {
		for v := range c {
			if !yield(v) {
				return
			}
		}
	}
}
