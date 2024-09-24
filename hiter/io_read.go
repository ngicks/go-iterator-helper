package hiter

import (
	"io"
	"iter"
)

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
