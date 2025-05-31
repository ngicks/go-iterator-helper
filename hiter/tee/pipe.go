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

// Pipe is similar [io.Pipe]
type Pipe[V any] struct {
	c      chan V
	closed chan struct{}
}

// NewPipe creates new Pipe instance.
// Unlike [io.Pipe], you can choose buffer size of internal channel by putting n.
// If n is less than or equal to zero, it is unbuffered.
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

// Close closes p.
// After calling Close, all Push failes returnning false,
// the iterator returned from [*Pipe.IntoIter] stops after buffered values.
//
// Close itself is not gorotine safe;
// it is safe to call Close adn Push from multiple goroutines.
// But simultaneous multiple calls to Close may panic.
func (p *Pipe[V]) Close() {
	select {
	case <-p.closed:
		return
	default:
	}
	close(p.closed)
	close(p.c)
}

// Push pushes v to p.
// Values might queue up in buffered channel if the internal channel was buffered,
// and might be received by the other end of the pipe
// by consuming an iterator retrived from [*Pipe.IntoIter].
//
// Push may block long if values are not read from the other side of the pipe.
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

// TryPush is like [*Pipe.Push], but does not block on sending.
// If [*Pipe.Close] is already called, open will be false.
// If it would block on sending v, pushed will false.
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

// IntoIter returns the reading side of pipe.
// The iterator yields values sent by [*Pipe.Push].
// The iterator stops only after [*Pipe.Close] is called.
func (p *Pipe[V]) IntoIter() iter.Seq[V] {
	return hiter.Chan(context.Background(), p.c)
}

// Pipes is like [Pipe], but it pipes key-value pairs.
type Pipe2[K, V any] struct {
	c      chan hiter.KeyValue[K, V]
	closed chan struct{}
}

// NewPipe2 creates new Pipe2 instance.
// Unlike [io.Pipe], you can choose buffer size of internal channel by putting n.
// If n is less than or equal to zero, it is unbuffered.
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

// Close closes p.
// After calling Close, all Push fails returnning false,
// the iterator returned from [*Pipe.IntoIter] stops after buffered values.
//
// Close itself is not gorotine safe;
// it is safe to call Close adn Push from multiple goroutines.
// But simultaneous multiple calls to Close may panic.
func (p *Pipe2[K, V]) Close() {
	select {
	case <-p.closed:
		return
	default:
	}
	close(p.closed)
	close(p.c)
}

// Push pushes a pair of k and v to p.
// Pairs might queue up in buffered channel if the internal channel was buffered,
// and will be received by the other end of the pipe
// by consuming an iterator retrived from [*Pipe.IntoIter].
//
// Push may block long if values are not read from the other side of the pipe.
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

// TryPush is like [*Pipe2.Push], but does not block on sending.
// If [*Pipe.Close] is already called, open will be false.
// If it would block on sending k adn v, pushed will false.
func (p *Pipe2[K, V]) TryPush(k K, v V) (open, pushed bool) {
	select {
	case <-p.closed:
		return false, false
	default:
	}
	select {
	case p.c <- hiter.KVPair(k, v):
		return true, true
	default:
		return true, false
	}
}

// IntoIter2 returns the reading side of pipe.
// The iterator yields k-v pairs sent by [*Pipe.Push].
// The iterator stops only after [*Pipe.Close] is called.
func (p *Pipe2[K, V]) IntoIter2() iter.Seq2[K, V] {
	return hiter.FromKeyValue(hiter.Chan(context.Background(), p.c))
}
