package iterable

import (
	"iter"
	"slices"

	"github.com/ngicks/go-iterator-helper/hiter"
)

var (
	_ hiter.IntoIterable[any]       = (*Peekable[any])(nil)
	_ hiter.IntoIterable2[any, any] = (*Peekable2[any, any])(nil)
)

// Peekable adds the read-ahead ability to [iter.Seq][V].
//
// The zero value of Peekable is not valid. Allocate one by [NewPeekable].
type Peekable[V any] struct {
	r      *Resumable[V]
	peeked []V
}

// NewPeekable initializes a peekable iterator.
// The caller must call [*Peekable.Stop] to release resources regardless of usage.
func NewPeekable[V any](seq iter.Seq[V]) *Peekable[V] {
	return &Peekable[V]{
		r: NewResumable(seq),
	}
}

// Stop releases resources allocated by [NewPeekable].
func (p *Peekable[V]) Stop() {
	p.r.Stop()
}

// Peek reads the next n elements without advancing the iterator.
// Peeked elements are only removed through the iterator returned from IntoIter.
func (p *Peekable[V]) Peek(n int) []V {
	// internal behavior might need some change to use ring buffer.
	if diff := n - len(p.peeked); diff > 0 {
		p.peeked = slices.AppendSeq(p.peeked, hiter.Limit(diff, p.r.IntoIter()))
	}
	if len(p.peeked) > n {
		return p.peeked[:n:n]
	}
	return slices.Clip(p.peeked)
}

func (p *Peekable[V]) pop() V {
	v0 := p.peeked[0]
	p.peeked[0] = *new(V)
	p.peeked = p.peeked[1:]
	return v0
}

// IntoIter returns p as an iterator form.
func (p *Peekable[V]) IntoIter() iter.Seq[V] {
	return func(yield func(V) bool) {
		for len(p.peeked) > 0 {
			if !yield(p.pop()) {
				return
			}
		}
		for v := range p.r.IntoIter() {
			if !yield(v) {
				return
			}
			for len(p.peeked) > 0 {
				if !yield(p.pop()) {
					return
				}
			}
		}
	}
}

// Peekable2 adds the read-ahead ability to [iter.Seq2][K, V].
//
// The zero value of Peekable2 is not valid. Allocate one by [NewPeekable2].
type Peekable2[K, V any] struct {
	r      *Resumable2[K, V]
	peeked []hiter.KeyValue[K, V]
}

// NewPeekable2 initializes a peekable iterator.
// The caller must call [*Peekable2.Stop] to release resources regardless of usage.
func NewPeekable2[K, V any](seq iter.Seq2[K, V]) *Peekable2[K, V] {
	return &Peekable2[K, V]{
		r: NewResumable2(seq),
	}
}

// Stop releases resources allocated by [NewPeekable2].
func (p *Peekable2[K, V]) Stop() {
	p.r.Stop()
}

// Peek reads the next n key-value pairs without advancing the iterator.
// Peeked pairs are only removed through the iterator returned from IntoIter.
func (p *Peekable2[K, V]) Peek(n int) []hiter.KeyValue[K, V] {
	if diff := n - len(p.peeked); diff > 0 {
		p.peeked = hiter.AppendSeq2(p.peeked, hiter.Limit2(diff, p.r.IntoIter2()))
	}
	if len(p.peeked) > n {
		return p.peeked[:n:n]
	}
	return slices.Clip(p.peeked)
}

func (p *Peekable2[K, V]) pop() (K, V) {
	kv0 := p.peeked[0]
	p.peeked[0] = *new(hiter.KeyValue[K, V])
	p.peeked = p.peeked[1:]
	return kv0.K, kv0.V
}

// IntoIter2 returns p as an iterator form.
func (p *Peekable2[K, V]) IntoIter2() iter.Seq2[K, V] {
	return func(yield func(K, V) bool) {
		for len(p.peeked) > 0 {
			if !yield(p.pop()) {
				return
			}
		}
		for k, v := range p.r.IntoIter2() {
			if !yield(k, v) {
				return
			}
			for len(p.peeked) > 0 {
				if !yield(p.pop()) {
					return
				}
			}
		}
	}
}
