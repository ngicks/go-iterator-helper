package goiteratorhelper_test

import (
	"fmt"
	"slices"

	"github.com/ngicks/go-iterator-helper/hiter"
	"github.com/ngicks/go-iterator-helper/hiter/iterable"
)

// Example_resumable demonstrates pausable and resumable iteration
func Example_resumable() {
	// Create a resumable iterator from a range
	source := hiter.Range(1, 10)
	resumable := iterable.NewResumable(source)

	// Process in batches
	fmt.Println("Processing in batches:")

	// Batch 1: Take first 3 elements
	batch1 := hiter.Limit(3, resumable.IntoIter())
	fmt.Println("Batch 1:", slices.Collect(batch1))

	// Batch 2: Take next 3 elements
	batch2 := hiter.Limit(3, resumable.IntoIter())
	fmt.Println("Batch 2:", slices.Collect(batch2))

	// Batch 3: Get remaining elements
	remaining := resumable.IntoIter()
	fmt.Println("Remaining:", slices.Collect(remaining))

	// Output:
	// Processing in batches:
	// Batch 1: [1 2 3]
	// Batch 2: [4 5 6]
	// Remaining: [7 8 9]
}