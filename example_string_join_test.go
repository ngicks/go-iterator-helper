package goiteratorhelper_test

import (
	"fmt"
	"slices"

	"github.com/ngicks/go-iterator-helper/hiter/stringsiter"
)

// Example_stringJoin demonstrates joining strings with an iterator
func Example_stringJoin() {
	words := slices.Values([]string{"go", "iterator", "helper"})
	collected := stringsiter.Collect(words)
	fmt.Println("Collected:", collected)
	joined := stringsiter.Join("-", words)
	fmt.Println("Joined:", joined)

	// Join with different separator
	parts := slices.Values([]string{"2024", "12", "25"})
	date := stringsiter.Join("/", parts)
	fmt.Println("Date:", date)

	// Output:
	// Collected: goiteratorhelper
	// Joined: go-iterator-helper
	// Date: 2024/12/25
}

