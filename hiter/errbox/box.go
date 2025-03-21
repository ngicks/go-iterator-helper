// errbox boxes iter.Seq[V, error] and converts to iter.Seq[V]. The occurrence of the error stops the boxed iterator. The error can be later inspected through method.
package errbox

import (
	"iter"

	"github.com/ngicks/go-iterator-helper/hiter/iterable"
)

// Box boxes an input [iter.Seq2][V, error] to [iter.Seq][V] by stripping nil errors from the input iterator.
// The first non-nil error causes the boxed iterator to be stopped and Box stores the error.
// Later the error can be examined by [*Box.Err].
//
// Box converts input iterator stateful form by calling [iter.Pull2].
// [*Box.IntoIter] returns the iterator as [iter.Seq][V].
// While consuming values from the iterator, it might conditionally yield a non-nil error.
// In that case Box stores the error and stops without yielding the value V paired to the error.
// [*Box.Err] returns that error otherwise nil.
//
// The zero Box is invalid and it must be allocated by [New].
type Box[V any] struct {
	err       error
	resumable *iterable.Resumable2[V, error]
}

// New returns a newly allocated Box.
//
// When a pair from seq contains non-nil error, Box discards a former value of that pair(V),
// then the iterator returned from [Box.IntoIter] stops.
//
// [*Box.Stop] must be called to release resource regardless of usage.
func New[V any](seq iter.Seq2[V, error]) *Box[V] {
	return &Box[V]{
		resumable: iterable.NewResumable2(seq),
	}
}

// Stop releases resources allocated by [New].
// After calling Stop, iterators returned from [Box.IntoIter] produce no more data.
func (b *Box[V]) Stop() {
	b.resumable.Stop()
}

// IntoIter returns an iterator which yields values from the input iterator.
//
// As the name IntoIter suggests, the iterator is stateful;
// breaking and calling seq again continues its iteration without replaying data.
// If the iterator finds an error, it stops iteration and will no longer produce any data.
// In that case the error can be inspected by calling [*Box.Err].
func (b *Box[V]) IntoIter() iter.Seq[V] {
	return func(yield func(V) bool) {
		for v, err := range b.resumable.IntoIter2() {
			if err != nil {
				b.Stop()
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
