package hiter

import (
	"io"
	"iter"
)

// Decode returns an iterator over consecutive decode results of dec.
//
// The iterator stops if and only if dec returns io.EOF. Handling other errors is caller's responsibility.
// If the first error should stop the iterator, use [LimitUntil], [LimitAfter] or [*errbox.Box].
func Decode[V any, Dec interface{ Decode(any) error }](dec Dec) iter.Seq2[V, error] {
	return func(yield func(V, error) bool) {
		for {
			var v V
			err := dec.Decode(&v)
			if err == io.EOF {
				return
			}
			if !yield(v, err) {
				return
			}
		}
	}
}
