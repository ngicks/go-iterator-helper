package iterable

import (
	"context"
	"iter"
	"sync"

	"github.com/ngicks/go-iterator-helper/hiter"
)

func TeeSeq[V any](seq iter.Seq[V], pusher func(v V) bool) iter.Seq[V] {
	return func(yield func(V) bool) {
		for v := range seq {
			if !pusher(v) {
				return
			}
			if !yield(v) {
				return
			}
		}
	}
}

func TeeSeqPipe[V any](seq iter.Seq[V]) (*Pipe[V], *Resumable[V]) {
	p := NewPipe[V](0)
	tee := NewResumable(
		hiter.TapLast(
			p.Close,
			TeeSeq(seq, p.Push),
		),
	)
	return p, tee
}

var (
	_ hiter.IntoIterable[any] = (*Pipe[any])(nil)
)

type Pipe[V any] struct {
	c         chan V
	closeOnce sync.Once
	closed    chan struct{}
}

func NewPipe[V any](n int) *Pipe[V] {
	if n < 0 {
		// panic?
		n = 0
	}
	p := &Pipe[V]{
		c:      make(chan V, n),
		closed: make(chan struct{}),
	}
	return p
}

func (p *Pipe[V]) Close() {
	p.closeOnce.Do(func() {
		close(p.closed)
		close(p.c)
	})
}

func (p *Pipe[V]) Push(v V) bool {
	select {
	case <-p.closed:
		return false
	default:
	}
	select {
	case p.c <- v:
		return true
	case <-p.closed:
		return false
	}
}

func (p *Pipe[V]) TryPush(v V) (ok, pushed bool) {
	select {
	case <-p.closed:
		return false, false
	default:
	}
	select {
	case p.c <- v:
		return true, true
	default:
		return true, false
	}
}

func (p *Pipe[V]) IntoIter() iter.Seq[V] {
	return hiter.Chan(context.Background(), p.c)
}
