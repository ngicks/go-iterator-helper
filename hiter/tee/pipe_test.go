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

// TestPipe2Operations tests Pipe2 specific operations
func TestPipe2Operations(t *testing.T) {
	t.Run("basic operations", func(t *testing.T) {
		p := NewPipe2[string, int](2) // Increase buffer to 2 to avoid blocking
		defer p.Close()

		// Test Push
		assert.Assert(t, p.Push("key1", 1))
		assert.Assert(t, p.Push("key2", 2))

		// Read values
		results := hiter.Collect2(hiter.Limit2(2, p.IntoIter2()))
		expected := []hiter.KeyValue[string, int]{
			{K: "key1", V: 1},
			{K: "key2", V: 2},
		}
		assert.DeepEqual(t, expected, results)
	})

	t.Run("TryPush operations", func(t *testing.T) {
		p := NewPipe2[string, int](0) // unbuffered

		// TryPush on empty unbuffered channel should fail
		open, pushed := p.TryPush("key", 1)
		assert.Assert(t, open)    // pipe is open
		assert.Assert(t, !pushed) // but push failed (would block)

		p.Close()

		// TryPush after close
		open, pushed = p.TryPush("key", 1)
		assert.Assert(t, !open)   // pipe is closed
		assert.Assert(t, !pushed) // push failed
	})

	t.Run("TryPush with buffer", func(t *testing.T) {
		p := NewPipe2[string, int](2)
		defer p.Close()

		// TryPush should succeed while buffer has space
		open, pushed := p.TryPush("key1", 1)
		assert.Assert(t, open && pushed)

		open, pushed = p.TryPush("key2", 2)
		assert.Assert(t, open && pushed)

		// Buffer is full, should fail
		open, pushed = p.TryPush("key3", 3)
		assert.Assert(t, open)    // pipe still open
		assert.Assert(t, !pushed) // but can't push (buffer full)

		// Read one value to make space
		k, v, ok := hiter.First2(p.IntoIter2())
		assert.Assert(t, ok)
		assert.Equal(t, "key1", k)
		assert.Equal(t, 1, v)

		// Now TryPush should succeed again
		open, pushed = p.TryPush("key3", 3)
		assert.Assert(t, open && pushed)
	})

	t.Run("close operations", func(t *testing.T) {
		p := NewPipe2[string, int](1)

		// Push a value
		assert.Assert(t, p.Push("key", 1))
		p.Close()

		// Push after close should fail
		assert.Assert(t, !p.Push("key2", 2))

		// Should still be able to read buffered value
		results := hiter.Collect2(p.IntoIter2())
		expected := []hiter.KeyValue[string, int]{{K: "key", V: 1}}
		assert.DeepEqual(t, expected, results)

		// Multiple close calls should not panic
		p.Close()
		p.Close()
	})
}

// TestNewPipe2NegativeBuffer tests NewPipe2 with negative buffer size
func TestNewPipe2NegativeBuffer(t *testing.T) {
	p := NewPipe2[string, int](-5) // Should be treated as 0
	defer p.Close()

	// Should behave like unbuffered
	open, pushed := p.TryPush("key", 1)
	assert.Assert(t, open && !pushed) // Open but can't push (would block)
}
