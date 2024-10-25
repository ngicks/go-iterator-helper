package sh

import (
	"crypto/rand"
	"io"
	"iter"
	mathRand "math/rand/v2"

	"github.com/ngicks/go-iterator-helper/hiter"
)

// why math/rand/v2 not exporting this type constraint.

type intType interface {
	~int | ~int8 | ~int16 | ~int32 | ~int64 |
		~uint | ~uint8 | ~uint16 | ~uint32 | ~uint64 | ~uintptr
}

// Rng returns an iterator over an infinite sequence of pseudo-random numbers in the half-open interval [0,n).
func Rng[Num intType](n Num) iter.Seq[Num] {
	return hiter.RepeatFunc(func() Num { return mathRand.N(n) }, -1)
}

// RngSourced is like [Rng] but accepts any arbitrary [mathRand.Source] implementations as a pseudo-random number source.
func RngSourced[Num intType](n Num, src mathRand.Source) iter.Seq[Num] {
	rng := mathRand.New(src)
	return hiter.RepeatFunc(func() Num { return Num(rng.Uint64N(uint64(n))) }, -1)
}

// RandBytes returns an iterator over pseudo-random n bytes long slice.
// The buffer which the iterator returns is reused and should not be retained by the callers.
// Callers should explicitly clone the slice by [Clone] if needed.
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
