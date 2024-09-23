package iterable

import (
	"iter"

	"github.com/ngicks/go-iterator-helper/hiter"
)

var (
	_ hiter.IntoIterable[any]       = (*Resumable[any])(nil)
	_ hiter.IntoIterable2[any, any] = (*Resumable2[any, any])(nil)
)

// Resumable converts the input [iter.Seq][V] into stateful form by calling [iter.Pull].
//
// The zero value of Resumable is not valid. Allocate one by [NewResumable].
type Resumable[V any] struct {
	next func() (V, bool)
	stop func()
}

// NewResumable wraps seq into stateful form so that the iterator can be break-ed and resumed.
// The caller must call [*Resumable.Stop] to release resources regardless of usage.
func NewResumable[V any](seq iter.Seq[V]) *Resumable[V] {
	next, stop := iter.Pull(seq)
	return &Resumable[V]{
		next: next,
		stop: stop,
	}
}

// Stop releases resources allocated by [NewResumable].
func (r *Resumable[V]) Stop() {
	r.stop()
}

// IntoIter returns an iterator over the input iterator.
// The iterator can be paused by break and later resumed without replaying data.
func (r *Resumable[V]) IntoIter() iter.Seq[V] {
	return func(yield func(V) bool) {
		for {
			v, ok := r.next()
			if !ok {
				return
			}
			if !yield(v) {
				return
			}
		}
	}
}

// Resumable2 converts the input [iter.Seq2][K, V] into stateful form by calling [iter.Pull2].
//
// The zero value of Resumable2 is not valid. Allocate one by [NewResumable2].
type Resumable2[K, V any] struct {
	next func() (K, V, bool)
	stop func()
}

// NewResumable2 wraps seq into stateful form so that the iterator can be break-ed and resumed.
// The caller must call [*Resumable2.Stop] to release resources regardless of usage.
func NewResumable2[K, V any](seq iter.Seq2[K, V]) *Resumable2[K, V] {
	next, stop := iter.Pull2(seq)
	return &Resumable2[K, V]{
		next: next,
		stop: stop,
	}
}

// Stop releases resources allocated by [NewResumable2].
func (r *Resumable2[K, V]) Stop() {
	r.stop()
}

// IntoIter2 returns an iterator over the input iterator.
// The iterator can be paused by break and later resumed without replaying data.
func (r *Resumable2[K, V]) IntoIter2() iter.Seq2[K, V] {
	return func(yield func(K, V) bool) {
		for {
			k, v, ok := r.next()
			if !ok {
				return
			}
			if !yield(k, v) {
				return
			}
		}
	}
}
