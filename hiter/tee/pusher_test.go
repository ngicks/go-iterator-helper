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

// TestMultiPusher2 tests the MultiPusher2 functionality
func TestMultiPusher2(t *testing.T) {
	t.Run("all pushers succeed", func(t *testing.T) {
		var results1, results2 []hiter.KeyValue[string, int]

		pusher1 := func(k string, v int) bool {
			results1 = append(results1, hiter.KeyValue[string, int]{K: k, V: v})
			return true
		}
		pusher2 := func(k string, v int) bool {
			results2 = append(results2, hiter.KeyValue[string, int]{K: k, V: v})
			return true
		}

		multiPusher := MultiPusher2(pusher1, pusher2)

		assert.Assert(t, multiPusher("a", 1))
		assert.Assert(t, multiPusher("b", 2))

		expected := []hiter.KeyValue[string, int]{
			{K: "a", V: 1},
			{K: "b", V: 2},
		}
		assert.DeepEqual(t, expected, results1)
		assert.DeepEqual(t, expected, results2)
	})

	t.Run("one pusher fails", func(t *testing.T) {
		var results1, results2 []hiter.KeyValue[string, int]

		pusher1 := func(k string, v int) bool {
			results1 = append(results1, hiter.KeyValue[string, int]{K: k, V: v})
			return true
		}
		pusher2 := func(k string, v int) bool {
			results2 = append(results2, hiter.KeyValue[string, int]{K: k, V: v})
			return false // This pusher always fails
		}

		multiPusher := MultiPusher2(pusher1, pusher2)

		// Should return false when any pusher fails
		assert.Assert(t, !multiPusher("a", 1))

		// First pusher should have been called, second should have been called and failed
		assert.DeepEqual(t, []hiter.KeyValue[string, int]{{K: "a", V: 1}}, results1)
		assert.DeepEqual(t, []hiter.KeyValue[string, int]{{K: "a", V: 1}}, results2)
	})

	t.Run("no pushers", func(t *testing.T) {
		multiPusher := MultiPusher2[string, int]()
		// Should return true when no pushers (vacuous truth)
		assert.Assert(t, multiPusher("a", 1))
	})
}

// TestPush2MoreCoverage tests the Push2 function with more coverage
func TestPush2MoreCoverage(t *testing.T) {
	t.Run("multiple pushers, first fails", func(t *testing.T) {
		input := []hiter.KeyValue[string, int]{
			{K: "a", V: 1},
			{K: "b", V: 2},
		}

		var results1, results2 []hiter.KeyValue[string, int]

		pusher1 := func(k string, v int) bool {
			results1 = append(results1, hiter.KeyValue[string, int]{K: k, V: v})
			return k != "b" // Fail on second item
		}
		pusher2 := func(k string, v int) bool {
			results2 = append(results2, hiter.KeyValue[string, int]{K: k, V: v})
			return true
		}

		result := Push2(hiter.Values2(input), pusher1, pusher2)
		assert.Assert(t, !result)

		// Both pushers see first value, but only pusher1 sees "b" before failing
		// When pusher1 fails, pusher2 is not called for "b"
		assert.DeepEqual(t, []hiter.KeyValue[string, int]{
			{K: "a", V: 1},
			{K: "b", V: 2},
		}, results1)
		assert.DeepEqual(t, []hiter.KeyValue[string, int]{
			{K: "a", V: 1},
		}, results2)
	})
}
