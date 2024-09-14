package errbox

import (
	"iter"
)

// Box boxes an input iter.Seq2[V, error] to iter.Seq[V],
// by storing a first non nil error encountered.
// Later the error can be examined by [*Box.Err].
//
// [*Box.Iter] returns an iterator converted to be stateful from input [iter.Seq2] by [iter.Pull2].
// When non nil error is yielded from the iterator,
// Box stores the error and the boxed iterator stops
// without yielding paired value V.
// [*Box.Err] returns that error otherwise nil.
//
// The zero Box is invalid and it must be allocated by [New].
type Box[V any] struct {
	err  error
	next func() (V, error, bool)
	stop func()
}

// New returns a newly allocate Box.
//
// When a pair from seq contains non-nil error, Box discards a former value of the pair(V),
// then the iterator returned from [Box.Iter] stops.
//
// [*Box.Stop] must be called to release resource regardless of usage.
func New[V any](seq iter.Seq2[V, error]) *Box[V] {
	next, stop := iter.Pull2(seq)
	return &Box[V]{
		next: next,
		stop: stop,
	}
}

// Stop releases resources allocated by [New].
// After calling Stop, iterators returned  from [Box.Iter] yields nothing.
func (b *Box[V]) Stop() {
	b.stop()
}

// Iter returns a stateful iterator which yields values from the input iterator.
func (b *Box[V]) Iter() iter.Seq[V] {
	return func(yield func(V) bool) {
		for {
			if b.err != nil {
				b.stop()
				return
			}
			v, err, ok := b.next()
			if !ok {
				return
			}
			if err != nil || !yield(v) {
				b.stop()
				b.err = err
				return
			}
		}
	}
}

// Err returns an error the input iterator has returned.
// If the iterator has not yet encountered an error, Err returns nil.
func (b *Box[V]) Err() error {
	return b.err
}
