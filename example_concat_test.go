package goiteratorhelper_test

import (
	"fmt"
	"slices"

	"github.com/ngicks/go-iterator-helper/hiter"
)

// Example_concat demonstrates combining multiple iterators
func Example_concat() {
	first := slices.Values([]int{1, 2, 3})
	second := slices.Values([]int{10, 11})
	third := slices.Values([]int{20, 21, 22})

	combined := hiter.Concat(first, second, third)
	fmt.Println("Combined:", slices.Collect(combined))

	// Output:
	// Combined: [1 2 3 10 11 20 21 22]
}