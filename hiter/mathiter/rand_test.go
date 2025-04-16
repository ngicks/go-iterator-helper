package mathiter_test

import (
	"crypto/rand"
	"io"
	mathRand "math/rand/v2"
	"slices"
	"testing"

	"github.com/ngicks/go-iterator-helper/hiter"
	"github.com/ngicks/go-iterator-helper/hiter/mathiter"
	"gotest.tools/v3/assert"
)

func TestRng(t *testing.T) {
	for i := range 10 {
		i = i + 1
		assert.Assert(t, hiter.Every(func(j int) bool { return 0 <= j && j < i }, hiter.Limit(100, mathiter.Rng(i))))
	}
}

func TestRngSource(t *testing.T) {
	var seed [32]byte
	_, err := io.ReadFull(rand.Reader, seed[:])
	if err != nil {
		panic(err)
	}
	for i := range 10 {
		i = (i + 1) * 10
		r := mathRand.New(mathRand.NewChaCha8(seed))
		numsExpected := make([]int, 100)
		for idx := range 100 {
			numsExpected[idx] = int(r.Uint64N(uint64(i)))
		}
		rng := mathiter.RngSourced(i, mathRand.NewChaCha8(seed))
		numsActual := slices.Collect(hiter.Limit(100, rng))
		assert.DeepEqual(t, numsExpected, numsActual)
	}
}
