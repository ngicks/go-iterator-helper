// errbox boxes iter.Seq[V, error] and converts to iter.Seq[V]. The occurrence of the error stops the boxed iterator. The error can be later inspected through method.
package errbox

import (
	"iter"

	"github.com/ngicks/go-iterator-helper/hiter"
)

var (
	_ hiter.IntoIterable[any]       = (*Box[any])(nil)
	_ hiter.IntoIterable2[any, any] = (*Box2[any, any])(nil)
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

// Map maps values over the input iterator and put them in Box.
//
// See doc comment for [Box] for detailed caveats.
func Map[V1, V2 any](mapper func(v V1) (V2, error), seq iter.Seq[V1]) *Box[V2] {
	return New(hiter.Divide(mapper, seq))
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

// Box2 is like [Box] but with [iter.Seq2].
//
// See doc comment for [Box] for detail.
type Box2[K, V any] struct {
	err error
	seq iter.Seq2[hiter.KeyValue[K, V], error]
}

// New2 is like [New] but returns Box2.
func New2[K, V any](seq iter.Seq2[hiter.KeyValue[K, V], error]) *Box2[K, V] {
	return &Box2[K, V]{
		seq: seq,
	}
}

// Map2 maps pairs of values over the input iterator using mapper and put them in Box2.
//
// See doc comment for [Box] for detailed caveats.
func Map2[K1, V1, K2, V2 any](mapper func(K1, V1) (K2, V2, error), seq iter.Seq2[K1, V1]) *Box2[K2, V2] {
	return New2(
		hiter.Map2(
			func(k K1, v V1) (hiter.KeyValue[K2, V2], error) {
				k2, v2, err := mapper(k, v)
				return hiter.KeyValue[K2, V2]{K: k2, V: v2}, err
			},
			seq,
		),
	)
}

// IntoIter2 is like [Box.IntoIter] but is for Box2.
func (b *Box2[K, V]) IntoIter2() iter.Seq2[K, V] {
	return func(yield func(K, V) bool) {
		if b.err != nil {
			return
		}
		for kv, err := range b.seq {
			if err != nil {
				b.err = err
				return
			}
			if !yield(kv.K, kv.V) {
				return
			}
		}
	}
}

// Err returns an error the input iterator has returned.
// If the iterator has not yet encountered an error, Err returns nil.
func (b *Box2[K, V]) Err() error {
	return b.err
}
