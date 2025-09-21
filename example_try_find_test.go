package goiteratorhelper_test

import (
	"errors"
	"fmt"
	"slices"

	"github.com/ngicks/go-iterator-helper/hiter"
)

// Example_tryFind demonstrates error-aware searching
func Example_tryFind() {
	// Create an iterator with potential errors
	data := hiter.Pairs(
		slices.Values([]int{1, 2, 3, 4, 5}),
		hiter.Concat(
			hiter.Repeat(error(nil), 2),
			hiter.Once(errors.New("processing error")),
			hiter.Repeat(error(nil), 2),
		),
	)

	// TryFind stops on first error
	val, idx, err := hiter.TryFind(func(n int) bool { return n > 2 }, data)
	if err != nil {
		fmt.Printf("Error encountered at search: %v\n", err)
		fmt.Printf("Last value: %d, index: %d\n", val, idx)
	} else {
		fmt.Printf("Found: value=%d at index=%d\n", val, idx)
	}

	// Output:
	// Error encountered at search: processing error
	// Last value: 0, index: -1
}