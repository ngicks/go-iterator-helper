package hiter_test

import (
	"context"
	"fmt"
	"math/rand/v2"
	"sync"
	"time"

	"github.com/ngicks/go-iterator-helper/hiter"
	"github.com/ngicks/go-iterator-helper/x/exp/xiter"
)

func ExampleChanSend() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	var (
		in, out = make(chan string), make(chan string)
		wg      sync.WaitGroup
	)

	wg.Add(1)
	go func() {
		defer wg.Done()
		hiter.ChanSend(ctx, in, hiter.Repeat("hey", 3))
	}()

	for range 3 {
		wg.Add(1)
		go func() {
			defer wg.Done()
			// super duper heavy tasks
			_, _ = hiter.ChanSend(
				ctx,
				out,
				hiter.Tap(
					func(_ string) {
						// sleep for random duration.
						// Ensuring moderate distribution among workers(goroutines.)
						time.Sleep(rand.N[time.Duration](100))
					},
					xiter.Map(
						func(s string) string { return "✨" + s + "✨" },
						hiter.Chan(ctx, in),
					),
				),
			)
		}()
	}

	for i, decorated := range hiter.Enumerate(hiter.Chan(ctx, out)) {
		fmt.Printf("%s\n", decorated)
		if i == 2 {
			cancel()
		}
	}
	wg.Wait()
	// Output:
	// ✨hey✨
	// ✨hey✨
	// ✨hey✨
}
