package iterable

import (
	"context"
	"iter"

	"github.com/ngicks/go-iterator-helper/hiter"
)

var _ hiter.IntoIterable[any] = Chan[any]{}

// Chan adds IntoIter method to a receive only channel.
type Chan[V any] struct {
	Ctx context.Context
	C   <-chan V
}

func (c Chan[V]) IntoIter() iter.Seq[V] {
	return hiter.Chan(c.Ctx, c.C)
}
