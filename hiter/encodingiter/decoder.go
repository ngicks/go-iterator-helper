package encodingiter

import (
	"io"
	"iter"

	"github.com/ngicks/go-iterator-helper/hiter"
)

var (
	_ = hiter.LimitUntil[any]
	_ = hiter.LimitAfter[any]
)

// Decode returns an iterator over consecutive decode results of dec.
//
// The iterator stops if and only if dec returns io.EOF. Handling other errors is caller's responsibility.
// If the first error should stop the iterator, use [hiter.LimitUntil], [hiter.LimitAfter] or even *errbox.Box.
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

// There won't be counter parts for [WriteTextMarshaler]/[WriteBinaryMarshaler]
// since the implementation of [encoding.TextUnmarshaler]/[encoding.TextUnmarshaler] might not know about boundary between
// 2 values from input byte stream.
