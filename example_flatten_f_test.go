package goiteratorhelper_test

import (
	"fmt"
	"slices"

	"github.com/ngicks/go-iterator-helper/hiter"
)

// Example_flattenF demonstrates FlattenF for pair flattening
func Example_flattenF() {
	// Categories with items
	categories := [][]string{{"fruits", "vegetables"}, {"proteins"}}
	counts := []int{10, 5}

	// FlattenF expands the first element (categories) while keeping the second
	expanded := hiter.FlattenF(hiter.Pairs(slices.Values(categories), slices.Values(counts)))

	fmt.Println("Category-Count pairs:")
	for category, count := range expanded {
		fmt.Printf("  %s: %d items\n", category, count)
	}

	// Output:
	// Category-Count pairs:
	//   fruits: 10 items
	//   vegetables: 10 items
	//   proteins: 5 items
}