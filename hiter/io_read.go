package hiter

import (
	"io"
	"io/fs"
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

// Readdir returns an iterator over dirents of a file-like object.
// On [io.EOF] the iterator stops without yielding the error.
func Readdir[R interface {
	Readdir(n int) ([]fs.FileInfo, error)
}](r R) iter.Seq2[fs.FileInfo, error] {
	return func(yield func(fs.FileInfo, error) bool) {
		for {
			// 64 dirents in a batch.
			dirents, err := r.Readdir(1 << 6)
			for _, dirent := range dirents {
				if !yield(dirent, nil) {
					return
				}
			}
			if err != nil {
				if err == io.EOF {
					return
				}
				yield(nil, err)
				return
			}
		}
	}
}
