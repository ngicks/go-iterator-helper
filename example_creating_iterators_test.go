package goiteratorhelper_test

import (
	"fmt"
	"slices"

	"github.com/ngicks/go-iterator-helper/hiter"
)

// Example_creatingIterators demonstrates creating iterators from various sources
func Example_creatingIterators() {
	// From slices
	slice_iter := slices.Values([]int{1, 2, 3})
	fmt.Println("From slice:", slices.Collect(slice_iter))

	// Using Range
	range_iter := hiter.Range(10, 15)
	fmt.Println("Range:", slices.Collect(range_iter))

	// Using Repeat
	repeat_iter := hiter.Repeat("hi", 3)
	fmt.Println("Repeat:", slices.Collect(repeat_iter))

	// Using Once
	once_iter := hiter.Once(42)
	fmt.Println("Once:", slices.Collect(once_iter))

	// Output:
	// From slice: [1 2 3]
	// Range: [10 11 12 13 14]
	// Repeat: [hi hi hi]
	// Once: [42]
}