package tee

import (
	"context"
	"io"
	"iter"

	"github.com/ngicks/go-iterator-helper/hiter"
)

var _ = io.Pipe

var (
	_ hiter.IntoIterable[any]       = (*Pipe[any])(nil)
	_ hiter.IntoIterable2[any, any] = (*Pipe2[any, any])(nil)
)

// Pipe is like [io.Pipe], but conveys values over an iterator.
type Pipe[V any] struct {
	c      chan V
	closed chan struct{}
}

// Experimental: not tested and might be changed any time.
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
	select {
	case <-p.closed:
		return
	default:
	}
	close(p.closed)
	close(p.c)
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

func (p *Pipe[V]) TryPush(v V) (open, pushed bool) {
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

type Pipe2[K, V any] struct {
	c      chan hiter.KeyValue[K, V]
	closed chan struct{}
}

func NewPipe2[K, V any](n int) *Pipe2[K, V] {
	if n < 0 {
		// panic?
		n = 0
	}
	p := &Pipe2[K, V]{
		c:      make(chan hiter.KeyValue[K, V], n),
		closed: make(chan struct{}),
	}
	return p
}

func (p *Pipe2[K, V]) Close() {
	select {
	case <-p.closed:
		return
	default:
	}
	close(p.closed)
	close(p.c)
}

func (p *Pipe2[K, V]) Push(k K, v V) bool {
	select {
	case <-p.closed:
		return false
	default:
	}
	select {
	case p.c <- hiter.KeyValue[K, V]{K: k, V: v}:
		return true
	case <-p.closed:
		return false
	}
}

func (p *Pipe2[K, V]) TryPush(k K, v V) (open, pushed bool) {
	select {
	case <-p.closed:
		return false, false
	default:
	}
	select {
	case p.c <- hiter.KeyValue[K, V]{K: k, V: v}:
		return true, true
	default:
		return true, false
	}
}

func (p *Pipe2[K, V]) IntoIter2() iter.Seq2[K, V] {
	return hiter.FromKeyValue(hiter.Chan(context.Background(), p.c))
}
