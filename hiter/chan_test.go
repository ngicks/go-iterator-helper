package hiter_test

import (
	"context"
	"iter"
	"sync/atomic"
	"testing"

	"github.com/ngicks/go-iterator-helper/hiter"
	"github.com/ngicks/go-iterator-helper/hiter/iterable"
)

func TestChanAll(t *testing.T) {
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

	testCase1[int]{
		Seq: func() iter.Seq[int] {
			return hiter.Chan(context.Background(), chanAll())
		},
		Seqs: []func() iter.Seq[int]{
			func() iter.Seq[int] {
				return iterable.Chan[int]{C: chanAll()}.IntoIter()
			},
		},
		Expected: []int{5, 6, 7, 8, 9, 10},
		BreakAt:  3,
	}.Test(t, func() { (*cancelP.Load())() })

	testCase1[int]{
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
	}.Test(t)
}
