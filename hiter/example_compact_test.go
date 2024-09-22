package hiter_test

import (
	"fmt"
	"slices"

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

func ExampleCompactFunc2() {
	type example struct {
		Key  string
		Data string
	}
	for i, v := range hiter.CompactFunc2(
		func(i1 int, v1 example, i2 int, v2 example) bool { return v1.Key == v2.Key },
		hiter.Enumerate(slices.Values([]example{
			{"foo", "yay"}, {"foo", "nay"}, {"foo", "mah"},
			{"bar", "yay"},
			{"baz", "yay"}, {"baz", "nay"},
		})),
	) {
		fmt.Printf("%d: %v\n", i, v)
	}
	// Output:
	// 0: {foo yay}
	// 3: {bar yay}
	// 4: {baz yay}
}
