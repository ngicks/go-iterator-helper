package cryptoiter

import (
	"crypto/rand"
	"io"
	"iter"

	"github.com/ngicks/go-iterator-helper/hiter"
	"github.com/ngicks/go-iterator-helper/hiter/mapper"
)

var _ = mapper.Clone[[]any]

// RandBytes returns an iterator over pseudo-random n bytes long slice.
// The buffer which the iterator returns is reused and should not be retained by the callers.
// Callers should explicitly clone the slice by [mapper.Clone] if needed.
func RandBytes(size int) iter.Seq[[]byte] {
	buf := make([]byte, size)
	return hiter.RepeatFunc(
		func() []byte {
			_, err := io.ReadFull(rand.Reader, buf)
			if err != nil {
				panic(err)
			}
			return buf
		},
		-1,
	)
}
