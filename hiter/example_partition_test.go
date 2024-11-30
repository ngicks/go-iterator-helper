package hiter_test

import (
	"fmt"
	"slices"

	"github.com/ngicks/go-iterator-helper/hiter"
)

// Example_partition demonstrates how you replace https://pkg.go.dev/github.com/samber/lo#PartitionBy with hiter.
func Example_partition() {
	mm := hiter.ReduceGroup(
		func(accumulator []int, current int) []int { return append(accumulator, current) },
		[]int(nil),
		hiter.WithGroupId(
			func(i int) string {
				if i < 0 {
					return "negative"
				} else if i%2 == 0 {
					return "even"
				}
				return "odd"
			},
			slices.Values([]int{-2, -1, 0, 1, 2, 3, 4, 5}),
		),
	)

	fmt.Printf(
		"%#v\n",
		slices.Collect(
			hiter.OmitF(
				hiter.MapKeys(
					mm,
					hiter.Range(0, len(mm)),
				),
			),
		),
	)
	// Output:
	// [][]int{[]int{-2, -1}, []int{0, 2, 4}, []int{1, 3, 5}}
}
