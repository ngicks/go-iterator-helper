package goiteratorhelper_test

import (
	"fmt"
	"slices"

	"github.com/ngicks/go-iterator-helper/hiter"
)

// Example_reduce demonstrates aggregation with Reduce
func Example_reduce() {
	numbers := slices.Values([]int{1, 2, 3, 4, 5})

	// Sum using Reduce
	sum := hiter.Reduce(func(acc, val int) int { return acc + val }, 0, numbers)
	fmt.Println("Sum:", sum)

	// Product using Reduce
	numbers2 := slices.Values([]int{1, 2, 3, 4, 5})
	product := hiter.Reduce(func(acc, val int) int { return acc * val }, 1, numbers2)
	fmt.Println("Product:", product)

	// Output:
	// Sum: 15
	// Product: 120
}

