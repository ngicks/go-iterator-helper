package hiter_test

import (
	"fmt"

	"github.com/ngicks/go-iterator-helper/hiter"
)

func ExamplePermutations() {
	in := []int{1, 2, 3}
	for p := range hiter.Permutations(in) {
		fmt.Printf("%#v\n", p)
	}
	// Output:
	// []int{1, 2, 3}
	// []int{2, 1, 3}
	// []int{3, 1, 2}
	// []int{1, 3, 2}
	// []int{2, 3, 1}
	// []int{3, 2, 1}
}
