package async

import (
	"context"
	"iter"
	"sync"

	"github.com/ngicks/go-iterator-helper/hiter"
)

type resultLike[V any] hiter.KeyValue[V, error]

// Map maps values from seq asynchronously.
//
// Map applies mapper to values from seq in separate goroutines while keeping output order.
// When the order does not need to be kept, just send all values to workers through a channel using [hiter.ChanSend]
// and receive results via other channel using [hiter.Chan].
//
// queueLimit limits max amounts of possible simultaneous resource allocated.
// queueLimit must not be less than 1, otherwise Map panics.
// workerLimit limits max possible concurrent invocation of mapper.
// workerLimit is ok to be less than 1.
// When queueLimit > workerLimit, the total number of workers is only limited by queueLimit.
//
// Cancelling ctx may stop the returned iterator early.
// mapper should respect the ctx, otherwise it delays the iterator to return.
//
// If mapper panics Map panics with the first panic value.
//
// The counter part of Map for [iter.Seq2][K, V] does not exist since mapper is allowed to return error as second ret value.
// If you need to map [iter.Seq2][K, V], convert it into [iter.Seq][V] by [hiter.ToKeyValue].
func Map[V1, V2 any](
	ctx context.Context,
	queueLimit int,
	workerLimit int,
	mapper func(context.Context, V1) (V2, error),
	seq iter.Seq[V1],
) iter.Seq2[V2, error] {
	if queueLimit <= 0 {
		panic("queueLimit must be greater than 0")
	}
	return func(yield func(V2, error) bool) {
		ctx, cancel := context.WithCancel(ctx)
		defer cancel()

		var (
			reservations = make(chan chan resultLike[V2], queueLimit-1)
			workerSem    chan struct{}
			wg           sync.WaitGroup
			panicVal     any
			panicOnce    sync.Once
		)
		if workerLimit > 0 {
			workerSem = make(chan struct{}, workerLimit)
		}

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
				wg.Done()
				close(reservations)
			}()
			defer recordPanicOnce()
			for v := range seq {
				rsv := make(chan resultLike[V2], 1)

				select {
				case <-ctx.Done():
					return
				case reservations <- rsv:
				}

				if workerSem != nil {
					select {
					case <-ctx.Done():
						// close rsv in all paths
						close(rsv)
						return
					case workerSem <- struct{}{}:
					}
				}

				wg.Add(1)
				go func() {
					defer func() {
						close(rsv)
						if workerSem != nil {
							<-workerSem
						}
						wg.Done()
					}()
					defer recordPanicOnce()
					v2, err := mapper(ctx, v)
					rsv <- resultLike[V2]{K: v2, V: err}
				}()
			}
		}()

		defer func() {
			cancel()
			wg.Wait()
			if panicVal != nil {
				panic(panicVal)
			}
		}()
		// record panic in case simultaneously multiple sources have panicked.
		defer recordPanicOnce()
		for rsv := range reservations {
			result, ok := <-rsv
			if !ok {
				// TODO: ignore when ok is false and return?
				continue
			}
			if !yield(result.K, result.V) {
				return
			}
		}
	}
}
