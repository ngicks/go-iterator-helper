package ih_test

import (
	"context"
	"iter"
	"sync/atomic"
	"testing"

	"github.com/ngicks/go-iterator-helper/ih"
	"github.com/ngicks/go-iterator-helper/ih/iterable"
	"gotest.tools/v3/assert"
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

	var count atomic.Int64
	testCase1[int]{
		Seq: func() iter.Seq[int] {
			return ih.Chan(chanAll(), func() { count.Add(1) })
		},
		Seqs: []func() iter.Seq[int]{
			func() iter.Seq[int] {
				return iterable.Chan[int](chanAll()).IntoIter()
			},
		},
		Expected: []int{5, 6, 7, 8, 9, 10},
		BreakAt:  3,
	}.Test(t, func() { (*cancelP.Load())() })
	assert.Assert(t, count.Load() == 2)
}
