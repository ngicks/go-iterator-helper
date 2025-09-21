package goiteratorhelper_test

import (
	"fmt"
	"slices"

	"github.com/ngicks/go-iterator-helper/hiter"
)

// Example_flattenL demonstrates FlattenL for second element flattening
func Example_flattenL() {
	// Products with multiple prices (different stores)
	products := []string{"laptop", "phone"}
	prices := [][]float64{{999.99, 1099.99, 899.99}, {599.99, 649.99}}

	// FlattenL expands the second element (prices) while keeping the first
	expanded := hiter.FlattenL(hiter.Pairs(slices.Values(products), slices.Values(prices)))

	fmt.Println("Product-Price combinations:")
	for product, price := range expanded {
		fmt.Printf("  %s: $%.2f\n", product, price)
	}

	// Output:
	// Product-Price combinations:
	//   laptop: $999.99
	//   laptop: $1099.99
	//   laptop: $899.99
	//   phone: $599.99
	//   phone: $649.99
}