package mathiter

import (
	"iter"
	"math/rand/v2"

	"github.com/ngicks/go-iterator-helper/hiter"
)

// why math/rand/v2 not exporting this type constraint.

type intType interface {
	~int | ~int8 | ~int16 | ~int32 | ~int64 |
		~uint | ~uint8 | ~uint16 | ~uint32 | ~uint64 | ~uintptr
}

// Rng returns an iterator over an infinite sequence of pseudo-random numbers.
// Numbers are distributed in the half-open interval [0,n).
func Rng[Num intType](n Num) iter.Seq[Num] {
	return hiter.RepeatFunc(func() Num { return rand.N(n) }, -1)
}

// RngSourced is like [Rng] but accepts any arbitrary [rand.Source] implementations as a pseudo-random number source.
func RngSourced[Num intType](n Num, src rand.Source) iter.Seq[Num] {
	rng := rand.New(src)
	return hiter.RepeatFunc(func() Num { return Num(rng.Uint64N(uint64(n))) }, -1)
}
