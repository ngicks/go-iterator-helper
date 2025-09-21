package goiteratorhelper_test

import (
	"fmt"
	"slices"

	"github.com/ngicks/go-iterator-helper/hiter"
)

// Example_mapAndFilter demonstrates transforming iterators
func Example_mapAndFilter() {
	numbers := slices.Values([]int{1, 2, 3, 4, 5, 6})

	// Map to double values
	doubled := hiter.Map(func(n int) int { return n * 2 }, numbers)
	fmt.Println("Doubled:", slices.Collect(doubled))

	// Filter even numbers
	numbers2 := slices.Values([]int{1, 2, 3, 4, 5, 6})
	evens := hiter.Filter(func(n int) bool { return n%2 == 0 }, numbers2)
	fmt.Println("Evens:", slices.Collect(evens))

	// Chain Map and Filter
	numbers3 := slices.Values([]int{1, 2, 3, 4, 5, 6})
	doubled_evens := hiter.Filter(
		func(n int) bool { return n%4 == 0 },
		hiter.Map(func(n int) int { return n * 2 }, numbers3),
	)
	fmt.Println("Doubled evens:", slices.Collect(doubled_evens))

	// Output:
	// Doubled: [2 4 6 8 10 12]
	// Evens: [2 4 6]
	// Doubled evens: [4 8 12]
}