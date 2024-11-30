package cryptoiter_test

import (
	"testing"

	"github.com/ngicks/go-iterator-helper/hiter/cryptoiter"
	"github.com/ngicks/go-iterator-helper/x/exp/xiter"
	"gotest.tools/v3/assert"
	"gotest.tools/v3/assert/cmp"
)

func TestRandBytes(t *testing.T) {
	for bin := range xiter.Limit(cryptoiter.RandBytes(3), 10) {
		assert.Assert(t, cmp.Len(bin, 3))
	}

	for bin := range xiter.Limit(cryptoiter.RandBytes(10), 10) {
		assert.Assert(t, cmp.Len(bin, 10))
	}
}
