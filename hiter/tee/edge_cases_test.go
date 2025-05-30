package tee

import (
	"slices"
	"sync"
	"testing"
	"time"

	"github.com/ngicks/go-iterator-helper/hiter"
	"gotest.tools/v3/assert"
)

// TestTeeSeq_PusherReturnsFalse tests behavior when pusher returns false mid-iteration
func TestTeeSeq_PusherReturnsFalse(t *testing.T) {
	type testCase struct {
		name       string
		input      []int
		stopAfter  int
		expected   []int
		pushedVals []int
	}

	tests := []testCase{
		{
			name:       "stop at first",
			input:      []int{1, 2, 3, 4, 5},
			stopAfter:  0,
			expected:   nil,      // no values yielded since pusher returns false immediately
			pushedVals: []int{1}, // pusher sees first value before returning false
		},
		{
			name:       "stop in middle",
			input:      []int{1, 2, 3, 4, 5},
			stopAfter:  2,
			expected:   []int{1, 2},    // values yielded before pusher returns false
			pushedVals: []int{1, 2, 3}, // pusher sees one more value than yielded
		},
		{
			name:       "empty input",
			input:      []int{},
			stopAfter:  0,
			expected:   nil,
			pushedVals: nil,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			var pushedVals []int
			count := 0
			pusher := func(v int) bool {
				pushedVals = append(pushedVals, v)
				if count >= tc.stopAfter {
					return false
				}
				count++
				return true
			}

			teeSeq := TeeSeq(slices.Values(tc.input), pusher)
			actual := slices.Collect(teeSeq)

			assert.DeepEqual(t, tc.expected, actual)
			assert.DeepEqual(t, tc.pushedVals, pushedVals)
		})
	}
}

// TestTeeSeq_EarlyBreak tests behavior when consumer breaks early
func TestTeeSeq_EarlyBreak(t *testing.T) {
	input := []int{1, 2, 3, 4, 5}
	var pushedVals []int
	pusher := func(v int) bool {
		pushedVals = append(pushedVals, v)
		return true
	}

	teeSeq := TeeSeq(slices.Values(input), pusher)
	var actual []int
	for v := range teeSeq {
		actual = append(actual, v)
		if v == 2 {
			break
		}
	}

	assert.DeepEqual(t, []int{1, 2}, actual)
	assert.DeepEqual(t, []int{1, 2}, pushedVals) // pusher should only see values up to break
}

// TestTeeSeq_MultipleIterations tests using the same tee iterator multiple times
func TestTeeSeq_MultipleIterations(t *testing.T) {
	input := []int{1, 2, 3}
	var allPushedVals []int
	pusher := func(v int) bool {
		allPushedVals = append(allPushedVals, v)
		return true
	}

	teeSeq := TeeSeq(slices.Values(input), pusher)

	// First iteration
	first := slices.Collect(teeSeq)
	// Second iteration
	second := slices.Collect(teeSeq)

	assert.DeepEqual(t, input, first)
	assert.DeepEqual(t, input, second)
	// Pusher should be called for each iteration
	assert.DeepEqual(t, []int{1, 2, 3, 1, 2, 3}, allPushedVals)
}

// TestTeeSeq2_PusherReturnsFalse tests TeeSeq2 with pusher returning false
func TestTeeSeq2_PusherReturnsFalse(t *testing.T) {
	input := []hiter.KeyValue[string, int]{
		{K: "a", V: 1},
		{K: "b", V: 2},
		{K: "c", V: 3},
	}
	var pushedPairs []hiter.KeyValue[string, int]
	count := 0
	pusher := func(k string, v int) bool {
		pushedPairs = append(pushedPairs, hiter.KeyValue[string, int]{K: k, V: v})
		if count >= 1 { // stop after second pair
			return false
		}
		count++
		return true
	}

	teeSeq := TeeSeq2(hiter.Values2(input), pusher)
	actual := hiter.Collect2(teeSeq)

	expected := []hiter.KeyValue[string, int]{
		{K: "a", V: 1},
	} // Only first item since pusher returns false at count=1
	expectedPushed := []hiter.KeyValue[string, int]{
		{K: "a", V: 1},
		{K: "b", V: 2},
	} // pusher sees second value and returns false

	assert.DeepEqual(t, expected, actual)
	assert.DeepEqual(t, expectedPushed, pushedPairs)
}

// TestPipe_CloseOperations tests various close scenarios
func TestPipe_CloseOperations(t *testing.T) {
	t.Run("multiple close calls", func(t *testing.T) {
		p := NewPipe[int](1)
		p.Close()
		p.Close() // should not panic
		p.Close() // should not panic

		// Push after close should return false
		assert.Assert(t, !p.Push(1))
	})

	t.Run("close while pushing", func(t *testing.T) {
		p := NewPipe[int](0) // unbuffered
		var wg sync.WaitGroup
		pushed := make(chan bool, 1)

		wg.Add(1)
		go func() {
			defer wg.Done()
			// This push will block on unbuffered channel
			result := p.Push(42)
			pushed <- result
		}()

		// Give push goroutine time to start and block
		time.Sleep(10 * time.Millisecond)
		p.Close()

		wg.Wait()
		// Push should return false due to close
		assert.Assert(t, !<-pushed)
	})

	t.Run("close while reading", func(t *testing.T) {
		p := NewPipe[int](1)
		p.Push(42)
		p.Close()

		// Should still be able to read buffered value
		values := slices.Collect(p.IntoIter())
		assert.DeepEqual(t, []int{42}, values)
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

// TestTeeSeq_EmptySequence tests tee behavior with empty sequences
func TestTeeSeq_EmptySequence(t *testing.T) {
	var pusherCalled bool
	pusher := func(v int) bool {
		pusherCalled = true
		return true
	}

	teeSeq := TeeSeq(slices.Values([]int{}), pusher)
	result := slices.Collect(teeSeq)

	assert.DeepEqual(t, []int(nil), result)
	assert.Assert(t, !pusherCalled, "Pusher should not be called for empty sequence")
}

// TestTeeSeq_Ordering tests that tee preserves ordering
func TestTeeSeq_Ordering(t *testing.T) {
	input := make([]int, 20) // Reduce size for simpler test
	for i := range input {
		input[i] = i
	}

	var pushedValues []int
	pusher := func(v int) bool {
		pushedValues = append(pushedValues, v)
		return true
	}

	teeSeq := TeeSeq(slices.Values(input), pusher)
	result := slices.Collect(teeSeq)

	// Both pusher and result should preserve order
	for i, v := range result {
		assert.Equal(t, i, v, "Result ordering broken at index %d", i)
	}

	for i, v := range pushedValues {
		assert.Equal(t, i, v, "Pusher ordering broken at index %d", i)
	}

	assert.Equal(t, len(input), len(result))
	assert.Equal(t, len(input), len(pushedValues))
}

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

// TestTeeSeq_PanicHandling tests behavior when pusher panics
func TestTeeSeq_PanicHandling(t *testing.T) {
	t.Run("pusher panics", func(t *testing.T) {
		input := []int{1, 2, 3}
		pusher := func(v int) bool {
			if v == 2 {
				panic("test panic")
			}
			return true
		}

		teeSeq := TeeSeq(slices.Values(input), pusher)

		// Should panic when reaching value 2
		defer func() {
			r := recover()
			assert.Assert(t, r != nil, "Expected panic")
			assert.Equal(t, "test panic", r)
		}()

		slices.Collect(teeSeq)
		t.Fatal("Should have panicked")
	})
}
