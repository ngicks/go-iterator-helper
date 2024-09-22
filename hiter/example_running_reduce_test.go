package hiter_test

import (
	"fmt"
	"math"

	"github.com/ngicks/go-iterator-helper/hiter"
)

func ExampleRunningReduce() {
	for i, r := range hiter.Enumerate(
		hiter.RunningReduce(
			func(accum int64, next int64, idx int) int64 { return accum + next },
			int64(0),
			hiter.Range[int64](1, math.MaxInt64),
		),
	) {
		if i >= 5 {
			break
		}
		fmt.Printf("%d\n", r)
	}
	// Output:
	// 1
	// 3
	// 6
	// 10
	// 15
}
