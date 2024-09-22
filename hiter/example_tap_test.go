package hiter_test

import (
	"fmt"

	"github.com/ngicks/go-iterator-helper/hiter"
)

func ExampleTap() {
	for i := range hiter.Tap(
		func(i int) {
			fmt.Printf("observed: %d\n", i)
		},
		hiter.Range(1, 4),
	) {
		fmt.Printf("yielded:  %d\n", i)
	}
	// Output:
	// observed: 1
	// yielded:  1
	// observed: 2
	// yielded:  2
	// observed: 3
	// yielded:  3
}
