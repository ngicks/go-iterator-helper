package hiter_test

import (
	"fmt"

	"github.com/ngicks/go-iterator-helper/hiter"
	"github.com/ngicks/go-iterator-helper/x/exp/xiter"
)

func ExampleCompact() {
	m := xiter.Merge(
		xiter.Map(func(i int) int { return 2 * i }, hiter.Range(1, 11)),
		xiter.Map(func(i int) int { return 1 << i }, hiter.Range(1, 11)),
	)

	first := true
	for i := range hiter.Compact(m) {
		if !first {
			fmt.Printf(", ")
		}
		fmt.Printf("%d", i)
		first = false
	}
	fmt.Println()
	// Output:
	// 2, 4, 6, 8, 10, 12, 14, 16, 18, 20, 32, 64, 128, 256, 512, 1024
}
