package cryptoiter_test

import (
	"testing"

	"github.com/ngicks/go-iterator-helper/hiter"
	"github.com/ngicks/go-iterator-helper/hiter/cryptoiter"
	"gotest.tools/v3/assert"
	"gotest.tools/v3/assert/cmp"
)

func TestRandBytes(t *testing.T) {
	for bin := range hiter.Limit(10, cryptoiter.RandBytes(3)) {
		assert.Assert(t, cmp.Len(bin, 3))
	}

	for bin := range hiter.Limit(10, cryptoiter.RandBytes(10)) {
		assert.Assert(t, cmp.Len(bin, 10))
	}
}
