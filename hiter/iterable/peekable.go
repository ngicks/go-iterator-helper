package iterable

import (
	"iter"
	"slices"

	"github.com/ngicks/go-iterator-helper/hiter"
	"github.com/ngicks/go-iterator-helper/x/exp/xiter"
)

var (
	_ hiter.IntoIterable[any]       = (*Peekable[any])(nil)
	_ hiter.IntoIterable2[any, any] = (*Peekable2[any, any])(nil)
)

type Peekable[V any] struct {
	r      *Resumable[V]
	peeked []V
}

func NewPeekable[V any](seq iter.Seq[V]) *Peekable[V] {
	return &Peekable[V]{
		r: NewResumable(seq),
	}
}

func (p *Peekable[V]) Stop() {
	p.r.Stop()
}

func (p *Peekable[V]) Peek(n int) []V {
	// internal behavior might need some change to use ring buffer.
	if diff := n - len(p.peeked); diff > 0 {
		p.peeked = slices.AppendSeq(p.peeked, xiter.Limit(p.r.IntoIter(), diff))
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

type Peekable2[K, V any] struct {
	r      *Resumable2[K, V]
	peeked []hiter.KeyValue[K, V]
}

func NewPeekable2[K, V any](seq iter.Seq2[K, V]) *Peekable2[K, V] {
	return &Peekable2[K, V]{
		r: NewResumable2(seq),
	}
}

func (p *Peekable2[K, V]) Stop() {
	p.r.Stop()
}

func (p *Peekable2[K, V]) Peek(n int) []hiter.KeyValue[K, V] {
	if diff := n - len(p.peeked); diff > 0 {
		p.peeked = hiter.AppendSeq2(p.peeked, xiter.Limit2(p.r.IntoIter2(), diff))
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
