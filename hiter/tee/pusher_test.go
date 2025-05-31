package tee

import (
	"slices"
	"testing"

	"github.com/ngicks/go-iterator-helper/hiter"
	"gotest.tools/v3/assert"
)

// TestPush_EdgeCases tests Push function edge cases
func TestPush_EdgeCases(t *testing.T) {
	t.Run("empty sequence", func(t *testing.T) {
		var called bool
		pusher := func(v int) bool {
			called = true
			return true
		}

		result := Push(slices.Values([]int{}), pusher)
		assert.Assert(t, result)
		assert.Assert(t, !called) // pusher should not be called for empty sequence
	})

	t.Run("empty pusher list", func(t *testing.T) {
		result := Push(slices.Values([]int{1, 2, 3}))
		assert.Assert(t, result) // should succeed with no pushers
	})

	t.Run("pusher returns false mid-sequence", func(t *testing.T) {
		var values []int
		count := 0
		pusher := func(v int) bool {
			values = append(values, v)
			count++
			return count < 3 // return false after 3rd call
		}

		result := Push(slices.Values([]int{1, 2, 3, 4, 5}), pusher)
		assert.Assert(t, !result)
		assert.DeepEqual(t, []int{1, 2, 3}, values)
	})
}

// TestPush2_EdgeCases tests Push2 function edge cases
func TestPush2_EdgeCases(t *testing.T) {
	t.Run("empty sequence", func(t *testing.T) {
		var called bool
		pusher := func(k string, v int) bool {
			called = true
			return true
		}

		result := Push2(hiter.Values2([]hiter.KeyValue[string, int]{}), pusher)
		assert.Assert(t, result)
		assert.Assert(t, !called)
	})

	t.Run("empty pusher list", func(t *testing.T) {
		input := []hiter.KeyValue[string, int]{
			{K: "a", V: 1},
			{K: "b", V: 2},
		}
		result := Push2(hiter.Values2(input))
		assert.Assert(t, result)
	})
}

// TestMultiPusher_EdgeCases tests MultiPusher with various edge cases
func TestMultiPusher_EdgeCases(t *testing.T) {
	t.Run("empty pusher list", func(t *testing.T) {
		multiPusher := MultiPusher[int]()
		// Should return true when no pushers to call
		assert.Assert(t, multiPusher(42))
	})

	t.Run("mixed pusher results", func(t *testing.T) {
		var calls []int
		pusher1 := func(v int) bool {
			calls = append(calls, 1)
			return true
		}
		pusher2 := func(v int) bool {
			calls = append(calls, 2)
			return false // This should stop the chain
		}
		pusher3 := func(v int) bool {
			calls = append(calls, 3)
			return true
		}

		multiPusher := MultiPusher(pusher1, pusher2, pusher3)
		result := multiPusher(42)

		assert.Assert(t, !result)
		// pusher3 should not be called since pusher2 returned false
		assert.DeepEqual(t, []int{1, 2}, calls)
	})

	t.Run("all pushers succeed", func(t *testing.T) {
		var values []int
		pusher1 := func(v int) bool {
			values = append(values, v*2)
			return true
		}
		pusher2 := func(v int) bool {
			values = append(values, v*3)
			return true
		}

		multiPusher := MultiPusher(pusher1, pusher2)
		result := multiPusher(5)

		assert.Assert(t, result)
		assert.DeepEqual(t, []int{10, 15}, values)
	})
}
