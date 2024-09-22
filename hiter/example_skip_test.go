package hiter_test

import (
	"fmt"

	"github.com/ngicks/go-iterator-helper/hiter"
)

func ExampleSkipLast() {
	fmt.Print("Without SkipLast: ")
	first := true
	for i := range hiter.Range(0, 10) {
		if !first {
			fmt.Print(", ")
		}
		fmt.Printf("%d", i)
		first = false
	}
	fmt.Println()
	fmt.Print("With SkipLast:    ")

	first = true
	for i := range hiter.SkipLast(5, hiter.Range(0, 10)) {
		if !first {
			fmt.Print(", ")
		}
		fmt.Printf("%d", i)
		first = false
	}
	fmt.Println()
	// Output:
	// Without SkipLast: 0, 1, 2, 3, 4, 5, 6, 7, 8, 9
	// With SkipLast:    0, 1, 2, 3, 4
}
