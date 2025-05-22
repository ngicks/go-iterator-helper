// errbox boxes iter.Seq[V, error] and converts to iter.Seq[V]. The occurrence of the error stops the boxed iterator. The error can be later inspected through method.
package errbox

import (
	"iter"
)

// Box boxes an input [iter.Seq2][V, error] to [iter.Seq][V] by stripping nil errors from values over the input iterator.
// The first non-nil error causes the boxed iterator to be stopped and Box stores the error.
// Later the error can be examined by [*Box.Err].
//
// Caveats: Box remembers the first non-nil error but it DOES NOT mean that Box converts an input iterator to be stateful.
// In case users need [Box.IntoIter] to be break-ed and resumed, users must ensure that the input is stateful, or is ok to replay it.
// To convert iterators to be stateful, [github.com/ngicks/go-iterator-helper/hiter/iterable.NewResumable] is useful.
//
// [*Box.IntoIter] returns the iterator as [iter.Seq][V].
// While consuming values from the iterator, it might conditionally yield a non-nil error.
// In that case Box stores the error and stops without yielding the value paired to the error.
// [*Box.Err] returns that error otherwise nil.
// After the first non-nil error, [*Box.IntoIter] return an iterator that yields nothing.
//
// The zero Box is invalid and it must be allocated by [New].
type Box[V any] struct {
	err error
	seq iter.Seq2[V, error]
}

// New returns a newly allocated Box.
//
// When a pair from seq contains non-nil error, Box discards a former value of that pair(V),
// then the iterator returned from [Box.IntoIter] stops.
//
// See doc comment for [Box] for detailed caveats.
func New[V any](seq iter.Seq2[V, error]) *Box[V] {
	return &Box[V]{
		seq: seq,
	}
}

// IntoIter returns an iterator which yields values from the input iterator.
//
// As the name IntoIter suggests, the iterator is (partially) stateful;
// If the iterator produce a non-nil error, it stops iteration without yielding paired value(V)
// and will no longer produce any data.
// In that case the error can be inspected by calling [*Box.Err].
func (b *Box[V]) IntoIter() iter.Seq[V] {
	return func(yield func(V) bool) {
		if b.err != nil {
			return
		}
		for v, err := range b.seq {
			if err != nil {
				b.err = err
				return
			}
			if !yield(v) {
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
