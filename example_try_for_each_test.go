package goiteratorhelper_test

import (
	"errors"
	"fmt"
	"slices"

	"github.com/ngicks/go-iterator-helper/hiter"
)

// Example_tryForEach demonstrates error-aware iteration
func Example_tryForEach() {
	// Create an iterator with errors at specific positions
	data := hiter.Pairs(
		slices.Values([]string{"apple", "banana", "cherry", "date"}),
		hiter.Concat(
			hiter.Repeat(error(nil), 2),
			hiter.Once(errors.New("bad fruit")),
			hiter.Once(error(nil)),
		),
	)

	count := 0
	err := hiter.TryForEach(func(fruit string) {
		count++
		fmt.Printf("Processing: %s\n", fruit)
	}, data)

	if err != nil {
		fmt.Printf("Stopped after %d items due to: %v\n", count, err)
	} else {
		fmt.Printf("Processed all %d items successfully\n", count)
	}

	// Output:
	// Processing: apple
	// Processing: banana
	// Stopped after 2 items due to: bad fruit
}