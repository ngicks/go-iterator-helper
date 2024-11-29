package hiter_test

import (
	"context"
	"iter"
	"sync"
	"sync/atomic"
	"testing"

	"github.com/ngicks/go-iterator-helper/hiter"
	"github.com/ngicks/go-iterator-helper/hiter/internal/testcase"
	"github.com/ngicks/go-iterator-helper/hiter/iterable"
	"gotest.tools/v3/assert"
	"gotest.tools/v3/assert/cmp"
)

func TestChan(t *testing.T) {
	var (
		cancelP atomic.Pointer[context.CancelFunc]
	)

	chanAll := func() chan int {
		ctx, cancel := context.WithCancel(context.Background())
		cancelP.Store(&cancel)

		c := make(chan int)
		go func() {
			for i := 5; i <= 10; i++ {
				select {
				case <-ctx.Done():
				case c <- i:
				}
			}
			cancel()
			close(c)
		}()
		return c
	}

	testcase.TestCase1[int]{
		Seq: func() iter.Seq[int] {
			return hiter.Chan(context.Background(), chanAll())
		},
		Seqs: []func() iter.Seq[int]{
			func() iter.Seq[int] {
				return iterable.Chan[int]{Ctx: context.Background(), C: chanAll()}.IntoIter()
			},
		},
		Expected: []int{5, 6, 7, 8, 9, 10},
		BreakAt:  3,
		Stateful: true,
	}.Test(t, func(_, _ int) { (*cancelP.Load())() })

	testcase.TestCase1[int]{
		Seq: func() iter.Seq[int] {
			ctx, cancel := context.WithCancel(context.Background())
			var count int
			return hiter.Tap(
				func(_ int) {
					count++
					if count == 4 {
						cancel()
					}
				},
				hiter.Chan(ctx, chanAll()),
			)
		},
		Expected: []int{5, 6, 7, 8},
		BreakAt:  3,
		Stateful: true,
	}.Test(t)
}

func TestChanSend(t *testing.T) {
	t.Run("sent all", func(t *testing.T) {
		var r []int
		var wg sync.WaitGroup
		wg.Add(1)
		c := make(chan int)
		go func() {
			defer wg.Done()
			for v := range c {
				r = append(r, v) // 5
			}
		}()
		v, ok := hiter.ChanSend(context.Background(), c, hiter.Range(5, 10))
		close(c)
		wg.Wait()
		assert.Assert(t, ok)
		assert.Equal(t, 0, v)
		assert.Assert(t, cmp.DeepEqual([]int{5, 6, 7, 8, 9}, r))
	})

	t.Run("cancel", func(t *testing.T) {
		ctx, cancel := context.WithCancel(context.Background())
		var r []int
		var wg sync.WaitGroup
		wg.Add(1)
		c := make(chan int)
		go func() {
			defer wg.Done()
			r = append(r, <-c) // 5
			r = append(r, <-c) // 6
			r = append(r, <-c) // 7
			cancel()
		}()
		v, ok := hiter.ChanSend(ctx, c, hiter.Range(5, 10))
		close(c)
		wg.Wait()
		assert.Assert(t, !ok)
		assert.Assert(t, v == 8)
		assert.Assert(t, cmp.DeepEqual([]int{5, 6, 7}, r))
	})
}
