package async

import (
	"context"
	"iter"
	"sync"
	"time"

	"github.com/jonboulle/clockwork"
	"github.com/ngicks/go-iterator-helper/hiter/sh"
)

// Keeping this line to let linking work.
// hope dead code elimination works well on this.
var _ any = sh.Clone[[]any]

var (
	clock = clockwork.NewRealClock()
)

// Chunk returns an iterator over consecutive values of up to n elements from seq.
//
// The returned iterator reuses the buffer it yields. Apply [sh.Clone] if the caller needs to retain slices.
//
// Chunk may yield slices where 0 < len(s) <= n.
// Values may be shorter than n if timeout > 0 and the timeout duration passed since last time seq generated a value.
//
// Chunk panics if n is less than 1.
func Chunk[V any](timeout time.Duration, n int, seq iter.Seq[V]) iter.Seq[[]V] {
	if n <= 0 {
		panic("n cannot be less than 1")
	}
	return func(yield func([]V) bool) {
		var wg sync.WaitGroup

		ctx, cancel := context.WithCancel(context.Background())

		var (
			ch        = make(chan V)
			panicVal  any
			panicOnce sync.Once
		)

		recordPanicOnce := func() {
			rec := recover()
			if rec != nil {
				cancel()
				panicOnce.Do(func() {
					panicVal = rec
				})
			}
		}

		wg.Add(1)
		go func() {
			defer func() {
				close(ch)
				wg.Done()
			}()
			defer recordPanicOnce()
			for v := range seq {
				select {
				case <-ctx.Done():
					return
				case ch <- v:
				}
			}
		}()

		timer := clock.NewTimer(timeout)
		// reset when first value is yielded.
		// As of Go 1.23, this is really safe.
		timer.Stop()

		defer func() {
			cancel()
			wg.Wait()
			if panicVal != nil {
				panic(panicVal)
			}
		}()
		// record panic in case both seq and consumer have panicked.
		defer recordPanicOnce()

		buf := make([]V, n)
		idx := 0
		for {
			select {
			case v, ok := <-ch:
				if !ok {
					if idx > 0 {
						yield(buf[:idx:idx])
					}
					return
				}

				buf[idx] = v
				idx++
				if idx != n {
					if timeout > 0 {
						timer.Reset(timeout)
					}
				} else {
					if !yield(buf) {
						return
					}
					timer.Stop()
					idx = 0
				}
			case <-timer.Chan():
				if idx > 0 {
					if !yield(buf[:idx:idx]) {
						return
					}
					idx = 0
				}
			}
		}
	}
}
