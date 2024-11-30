package async

import (
	"context"
	"slices"
	"testing"
	"time"

	"github.com/jonboulle/clockwork"
	"github.com/ngicks/go-iterator-helper/hiter"
	"github.com/ngicks/go-iterator-helper/hiter/mapper"
	"github.com/ngicks/go-iterator-helper/x/exp/xiter"
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
		slices.Collect(mapper.Clone(xiter.Limit(Chunk(0, 3, src), 2))),
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
	defer func() {
		clock = org
	}()
	f := newResetCountClock(clockwork.NewFakeClock())
	clock = f

	resultChan := make(chan []int)
	go func() {
		for i, v := range hiter.Enumerate(Chunk(time.Millisecond, 3, src)) {
			t.Logf("%d: %#v", i, v)
			resultChan <- v
		}
		close(resultChan)
	}()

	waitReset := func() {
		t.Helper()
		assert.Equal(t, time.Millisecond, <-f.ResetChan())
	}

	c <- 0
	waitReset()
	c <- 1
	waitReset()
	f.BlockUntil(1)
	f.Advance(time.Millisecond + 100)
	assert.DeepEqual(t, []int{0, 1}, <-resultChan)
	c <- 2
	waitReset()
	c <- 3
	waitReset()
	c <- 4
	assert.DeepEqual(t, []int{2, 3, 4}, <-resultChan)
	c <- 5
	waitReset()
	c <- 6
	waitReset()
	f.BlockUntil(1)
	f.Advance(time.Millisecond)
	assert.DeepEqual(t, []int{5, 6}, <-resultChan)
	c <- 7
	waitReset()
	c <- 8
	waitReset()
	c <- 9
	assert.DeepEqual(t, []int{7, 8, 9}, <-resultChan)
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
