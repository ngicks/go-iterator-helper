package cryptoiter_test

import (
	"slices"
	"testing"

	"github.com/ngicks/go-iterator-helper/hiter"
	"github.com/ngicks/go-iterator-helper/hiter/cryptoiter"
	"github.com/ngicks/go-iterator-helper/hiter/mapper"
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

func TestRandBytesWithClone(t *testing.T) {
	// Test with mapper.Clone to ensure buffer reuse is covered
	cloned := slices.Collect(hiter.Limit(5, mapper.Clone(cryptoiter.RandBytes(4))))

	assert.Assert(t, cmp.Len(cloned, 5))
	for _, bin := range cloned {
		assert.Assert(t, cmp.Len(bin, 4))
	}

	// Verify buffers are different (not sharing memory)
	// Check if they point to different memory locations by modifying one
	original := make([]byte, len(cloned[0]))
	copy(original, cloned[0])
	cloned[0][0] = ^cloned[0][0] // flip bits
	assert.Assert(t, string(cloned[0]) != string(original))
}

func TestRandBytesZeroSize(t *testing.T) {
	// Test zero-size buffer
	for bin := range hiter.Limit(3, cryptoiter.RandBytes(0)) {
		assert.Assert(t, cmp.Len(bin, 0))
	}
}

func TestRandBytesLargeSize(t *testing.T) {
	// Test larger buffer size
	for bin := range hiter.Limit(2, cryptoiter.RandBytes(100)) {
		assert.Assert(t, cmp.Len(bin, 100))
	}
}
