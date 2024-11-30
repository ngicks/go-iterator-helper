package ioiter

import (
	"io"
	"io/fs"
	"iter"
)

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
