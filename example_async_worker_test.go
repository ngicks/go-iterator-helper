package goiteratorhelper_test

import (
	"context"
	"fmt"
	"maps"
	"slices"
	"sync"

	"github.com/ngicks/go-iterator-helper/hiter"
	"github.com/ngicks/go-iterator-helper/hiter/async"
	"github.com/ngicks/go-iterator-helper/hiter/iterable"
	"github.com/ngicks/go-iterator-helper/x/exp/xiter"
)

// Example async worker channel demonstrates usage of [hiter.Chan], [hiter.ChanSend].
// It sends values from seq to worker running on separates goroutines.
// Workers work on values and then send results back to the main goroutine.
func Example_async_worker_channel() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	works := []string{"foo", "bar", "baz"}

	in := make(chan string, 5)
	out := make(chan hiter.KeyValue[string, error])

	var wg sync.WaitGroup
	wg.Add(3)
	for range 3 {
		go func() {
			defer wg.Done()
			_, _ = hiter.ChanSend(
				ctx,
				out,
				xiter.Map(
					func(s string) hiter.KeyValue[string, error] {
						return hiter.KeyValue[string, error]{
							K: "✨" + s + "✨" + s + "✨",
							V: nil,
						}
					},
					hiter.Chan(ctx, in),
				),
			)
		}()
	}

	var wg2 sync.WaitGroup
	wg2.Add(1)
	go func() {
		defer wg2.Done()
		wg.Wait()
		close(out)
	}()

	_, _ = hiter.ChanSend(ctx, in, slices.Values(works))
	close(in)

	results := maps.Collect(hiter.FromKeyValue(hiter.Chan(ctx, out)))

	for result, err := range iterable.MapSorted[string, error](results).Iter2() {
		fmt.Printf("result = %s, err = %v\n", result, err)
	}

	wg2.Wait()

	// Output:
	// result = ✨bar✨bar✨, err = <nil>
	// result = ✨baz✨baz✨, err = <nil>
	// result = ✨foo✨foo✨, err = <nil>
}

// Example async worker map demonstrates usage of async.Map.
// At the surface it is similar to [xiter.Map2]. Actually it calls mapper in separate goroutine.
// If you don't care about order of element,
// just send values to workers through a channel and send back through another channel.
func Example_async_worker_map() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	works := []string{"foo", "bar", "baz"}

	// The order is kept.
	for result, err := range async.Map(
		ctx,
		/*queueLimit*/ 10,
		/*workerLimit*/ 5,
		/*mapper*/ func(ctx context.Context, s string) (string, error) {
			return "✨" + s + "✨" + s + "✨", nil
		},
		slices.Values(works),
	) {
		fmt.Printf("result = %s, err = %v\n", result, err)
	}
	// Output:
	// result = ✨foo✨foo✨, err = <nil>
	// result = ✨bar✨bar✨, err = <nil>
	// result = ✨baz✨baz✨, err = <nil>
}
