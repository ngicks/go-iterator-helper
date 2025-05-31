package tee

import (
	"testing"

	"github.com/ngicks/go-iterator-helper/hiter"
	"gotest.tools/v3/assert"
)

// TestPipe_BufferEdgeCases tests buffer capacity edge cases
func TestPipe_BufferEdgeCases(t *testing.T) {
	t.Run("exactly full buffer", func(t *testing.T) {
		p := NewPipe[int](2)
		defer p.Close()

		// Fill buffer exactly
		assert.Assert(t, p.Push(1))
		assert.Assert(t, p.Push(2))

		// Next push should block, test with TryPush
		open, pushed := p.TryPush(3)
		assert.Assert(t, open)    // pipe is still open
		assert.Assert(t, !pushed) // but push failed due to full buffer

		// Read one value to make space
		v, ok := hiter.First(p.IntoIter())
		assert.Assert(t, ok)
		assert.Equal(t, 1, v)

		// Now push should succeed
		assert.Assert(t, p.Push(3))
	})

	t.Run("negative buffer size", func(t *testing.T) {
		p := NewPipe[int](-5) // should be converted to 0
		defer p.Close()

		// Should behave like unbuffered channel
		open, pushed := p.TryPush(1)
		assert.Assert(t, open)
		assert.Assert(t, !pushed) // unbuffered channel blocks immediately
	})
}

// TestPipe_TryPushEdgeCases tests TryPush behavior
func TestPipe_TryPushEdgeCases(t *testing.T) {
	t.Run("try push to closed pipe", func(t *testing.T) {
		p := NewPipe[int](1)
		p.Close()

		open, pushed := p.TryPush(42)
		assert.Assert(t, !open)
		assert.Assert(t, !pushed)
	})

	t.Run("try push to full buffer", func(t *testing.T) {
		p := NewPipe[int](1)
		defer p.Close()

		// Fill buffer
		assert.Assert(t, p.Push(1))

		// Try push to full buffer
		open, pushed := p.TryPush(2)
		assert.Assert(t, open)    // pipe is open
		assert.Assert(t, !pushed) // but buffer is full
	})
}
