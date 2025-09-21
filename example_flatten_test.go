package goiteratorhelper_test

import (
	"fmt"
	"slices"

	"github.com/ngicks/go-iterator-helper/hiter"
)

// Example_flatten demonstrates flattening nested structures
func Example_flatten() {
	// Nested arrays representing grouped data
	grouped := [][]string{
		{"apple", "banana"},
		{"carrot"},
		{"dog", "elephant", "fox"},
		{"guitar"},
	}

	// Flatten the nested structure
	flattened := hiter.Flatten(slices.Values(grouped))

	fmt.Println("Grouped:", grouped)
	fmt.Println("Flattened:", slices.Collect(flattened))

	// Output:
	// Grouped: [[apple banana] [carrot] [dog elephant fox] [guitar]]
	// Flattened: [apple banana carrot dog elephant fox guitar]
}