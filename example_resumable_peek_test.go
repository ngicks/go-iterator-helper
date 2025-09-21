package goiteratorhelper_test

import (
	"fmt"

	"github.com/ngicks/go-iterator-helper/hiter"
	"github.com/ngicks/go-iterator-helper/hiter/iterable"
)

// Example_resumablePeek demonstrates peeking with resumable iterators
func Example_resumablePeek() {
	// Data stream to process
	data := hiter.Range(100, 106)
	resumable := iterable.NewResumable(data)

	// Peek at first element without consuming
	first, _ := hiter.First(resumable.IntoIter())
	fmt.Printf("Peeked first: %d\n", first)

	// Process all elements (including the peeked one)
	fmt.Print("All elements: ")
	for val := range resumable.IntoIter() {
		fmt.Printf("%d ", val)
	}
	fmt.Println()

	// Output:
	// Peeked first: 100
	// All elements: 101 102 103 104 105
}