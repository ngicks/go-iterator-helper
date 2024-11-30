package tee

import (
	"context"
	"io"
	"iter"

	"github.com/ngicks/go-iterator-helper/hiter"
	"github.com/ngicks/go-iterator-helper/hiter/iterable"
)

var _ = io.Reader(nil)

// TeeSeq is [iter.Seq] equivalent of [io.TeeReader].
//
// TeeSeq returns a [iter.Seq] that pushes to pusher what it reads from seq.
// Yielding values from the returned iterator performs push before the inner loop receives the value.
// The iterator is not stateful; you may want to wrap it with [iterable.Resumable].
// If pusher returns false, the iterator stops iteration without yielding value.
//
// Experimental: not tested and might be changed any time.
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

// TeeSeq2 is [iter.Seq2] equivalent of [io.TeeReader].
//
// TeeSeq2 returns a [iter.Seq2] that pushes to pusher what it reads from seq.
// Yielding key-value pairs from the returned iterator performs push before the inner loop receives the pair.
// The iterator is not stateful; you may want to wrap it with [iterable.Resumable2].
// If pusher returns false, the iterator stops iteration without yielding pair.
//
// Experimental: not tested and might be changed any time.
func TeeSeq2[K, V any](seq iter.Seq2[K, V], pusher func(K, V) bool) iter.Seq2[K, V] {
	return func(yield func(K, V) bool) {
		for k, v := range seq {
			if !pusher(k, v) {
				return
			}
			if !yield(k, v) {
				return
			}
		}
	}
}

// TeeSeqPipe tees values from seq to [*Pipe]. see doc comments for [TeeSeq].
// Yielding values from returned [*iterable.Resumable] also performs push to [*Pipe].
//
// Experimental: not tested and might be changed any time.
func TeeSeqPipe[V any](bufSize int, seq iter.Seq[V]) (*Pipe[V], *iterable.Resumable[V]) {
	p := NewPipe[V](bufSize)
	tee := iterable.NewResumable(
		hiter.TapLast(
			p.Close,
			TeeSeq(seq, p.Push),
		),
	)
	return p, tee
}

// TeeSeqPipe2 tees key-value pairs from seq to [*Pipe2]. see doc comments for [TeeSeq2].
// Yielding pairs from returned [*iterable.Resumable2] also performs push to [*Pipe2].
//
// Experimental: not tested and might be changed any time.
func TeeSeqPipe2[K, V any](bufSize int, seq iter.Seq2[K, V]) (*Pipe2[K, V], *iterable.Resumable2[K, V]) {
	p := NewPipe2[K, V](bufSize)
	tee := iterable.NewResumable2(
		hiter.TapLast2(
			p.Close,
			TeeSeq2(seq, p.Push),
		),
	)
	return p, tee
}

var (
	_ hiter.IntoIterable[any]       = (*Pipe[any])(nil)
	_ hiter.IntoIterable2[any, any] = (*Pipe2[any, any])(nil)
)

// Experimental: not tested and might be changed any time.
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

// Experimental: not tested and might be changed any time.
type Pipe2[K, V any] struct {
	c      chan hiter.KeyValue[K, V]
	closed chan struct{}
}

// Experimental: not tested and might be changed any time.
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
