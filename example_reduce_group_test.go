package goiteratorhelper_test

import (
	"fmt"
	"slices"

	"github.com/ngicks/go-iterator-helper/hiter"
)

// Example_reduceGroup demonstrates grouping and aggregation
func Example_reduceGroup() {
	// Group by key and sum values
	pairs := hiter.Pairs(
		slices.Values([]string{"fruit", "vegetable", "fruit", "fruit", "vegetable"}),
		slices.Values([]int{10, 20, 30, 15, 25}),
	)

	grouped := hiter.ReduceGroup(
		func(acc, val int) int { return acc + val },
		0,
		pairs,
	)
	fmt.Println("Grouped sums:", grouped)

	// Output:
	// Grouped sums: map[fruit:55 vegetable:45]
}