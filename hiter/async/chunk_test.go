package async

import (
	"context"
	"slices"
	"testing"
	"time"

	"github.com/jonboulle/clockwork"
	"github.com/ngicks/go-iterator-helper/hiter"
	"github.com/ngicks/go-iterator-helper/hiter/mapper"
	"gotest.tools/v3/assert"
	"gotest.tools/v3/assert/cmp"
)

type resetCountClock struct {
	resetCh chan time.Duration
	clockwork.FakeClock
}

func newResetCountClock(fake clockwork.FakeClock) *resetCountClock {
	return &resetCountClock{
		resetCh:   make(chan time.Duration),
		FakeClock: fake,
	}
}

func (r *resetCountClock) ResetChan() <-chan time.Duration {
	return r.resetCh
}

func (r *resetCountClock) NewTimer(d time.Duration) clockwork.Timer {
	return &resetCountTimer{
		resetCh: r.resetCh,
		Timer:   r.FakeClock.NewTimer(d),
	}
}

type resetCountTimer struct {
	resetCh chan time.Duration
	clockwork.Timer
}

func (r *resetCountTimer) Reset(d time.Duration) bool {
	r.resetCh <- d
	return r.Timer.Reset(d)
}

func TestChunk_successful(t *testing.T) {
	src := hiter.Range(0, 10)
	assert.DeepEqual(
		t,
		[][]int{{0, 1, 2}, {3, 4, 5}, {6, 7, 8}, {9}},
		slices.Collect(mapper.Clone(Chunk(0, 3, src))),
	)
	assert.DeepEqual(
		t,
		[][]int{{0, 1, 2}, {3, 4, 5}},
		slices.Collect(mapper.Clone(hiter.Limit(2, Chunk(0, 3, src)))),
	)
	assert.DeepEqual(
		t,
		[][]int(nil),
		slices.Collect(mapper.Clone(Chunk(0, 3, hiter.Empty[int]()))),
	)
}

func TestChunk_timeout(t *testing.T) {
	c := make(chan int)
	src := hiter.Chan(context.Background(), c)

	org := clock
	f := newResetCountClock(clockwork.NewFakeClock())
	clock = f
	defer func() {
		clock = org
	}()

	resultChan := make(chan []int)
	go func() {
		for i, v := range hiter.Enumerate(Chunk(time.Millisecond, 3, src)) {
			v = slices.Clone(v)
			t.Logf("%d: %#v", i, v)
			resultChan <- v
		}
		close(resultChan)
	}()

	waitReset := func() {
		t.Helper()
		assert.Equal(t, time.Millisecond, <-f.ResetChan())
	}

	// Test timeout with single value
	c <- 0
	waitReset()
	f.BlockUntil(1)
	f.Advance(time.Millisecond + 100)
	result1 := <-resultChan
	assert.Equal(t, 1, len(result1))
	assert.Equal(t, 0, result1[0])

	// Test that next value starts new chunk
	c <- 1
	waitReset()
	f.BlockUntil(1)
	f.Advance(time.Millisecond + 100)
	result2 := <-resultChan
	assert.Equal(t, 1, len(result2))
	assert.Equal(t, 1, result2[0])

	// Test filling buffer completely (no timeout)
	c <- 2
	waitReset()
	c <- 3
	c <- 4
	assert.DeepEqual(t, []int{2, 3, 4}, <-resultChan)

	f.Advance(time.Millisecond + 100)

	// Test timeout behavior - values sent close together may or may not be in same chunk
	c <- 5
	c <- 6
	waitReset() // This waits for timer reset from value 5
	f.BlockUntil(1)
	f.Advance(time.Millisecond + 100)
	result3 := <-resultChan
	// Can't strictly sync with this race.
	// TODO: add (var valueReceived chan struct{}) at top of module
	// and manipulate the channel to sync value received / time-out
	assert.Assert(t, len(result3) >= 1 && len(result3) <= 2)
	assert.Equal(t, 5, result3[0])
	if len(result3) == 2 {
		assert.Equal(t, 6, result3[1])
	}

	if len(result3) == 1 {
		waitReset()
		f.BlockUntil(1)
		f.Advance(time.Millisecond + 100)
		result4 := <-resultChan
		assert.Equal(t, 6, result4[0])
	}

	close(c)
	_, ok := <-resultChan
	assert.Assert(t, !ok)
}

func TestChunk_panic_propagation(t *testing.T) {
	src := hiter.Range(0, 10)
	t.Run("seq panics", func(t *testing.T) {
		var result [][]int
		func() {
			defer func() {
				rec := recover()
				assert.Assert(t, rec == errSample)
			}()
			for v := range mapper.Clone(
				Chunk(
					0,
					4,
					hiter.Tap(
						func(i int) {
							if i == 5 {
								panic(errSample)
							}
						},
						src,
					),
				),
			) {
				result = append(result, v)
			}
		}()
		assert.Assert(t, cmp.DeepEqual([][]int{{0, 1, 2, 3}, {4}}, result))
	})
	t.Run("consumer panics", func(t *testing.T) {
		var result [][]int
		func() {
			defer func() {
				rec := recover()
				assert.Assert(t, rec == errSample)
			}()
			for i, v := range hiter.Enumerate(mapper.Clone(Chunk(0, 4, src))) {
				result = append(result, v)
				if i == 1 {
					panic(errSample)
				}
			}
		}()
		assert.Assert(t, cmp.DeepEqual([][]int{{0, 1, 2, 3}, {4, 5, 6, 7}}, result))
	})
}
