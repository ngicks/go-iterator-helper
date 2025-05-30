package hiter_test

import (
	"fmt"
	"iter"
	"slices"

	"github.com/ngicks/go-iterator-helper/hiter"
)

func ExampleWindow_moving_average() {
	src := []int{1, 0, 1, 0, 1, 0, 5, 3, 2, 3, 4, 6, 5, 3, 6, 7, 7, 8, 9, 5, 7, 7, 8}
	first := true
	for avg := range hiter.Map(
		func(s []int) float64 {
			return float64(hiter.Sum(slices.Values(s))) / float64(len(s))
		},
		hiter.Window(src, 5),
	) {
		if !first {
			fmt.Print(", ")
		}
		fmt.Printf("%02.1f", avg)
		first = false
	}
	fmt.Println()
	// Output:
	// 0.6, 0.4, 1.4, 1.8, 2.2, 2.6, 3.4, 3.6, 4.0, 4.2, 4.8, 5.4, 5.6, 6.2, 7.4, 7.2, 7.2, 7.2, 7.2
}

func ExampleWindowSeq_moving_average() {
	src := []int{1, 0, 1, 0, 1, 0, 5, 3, 2, 3, 4, 6, 5, 3, 6, 7, 7, 8, 9, 5, 7, 7, 8}
	first := true
	for avg := range hiter.Map(
		func(s iter.Seq[int]) float64 {
			return float64(hiter.Sum(s)) / float64(5)
		},
		hiter.WindowSeq(5, slices.Values(src)),
	) {
		if !first {
			fmt.Print(", ")
		}
		fmt.Printf("%02.1f", avg)
		first = false
	}
	fmt.Println()
	// Output:
	// 0.6, 0.4, 1.4, 1.8, 2.2, 2.6, 3.4, 3.6, 4.0, 4.2, 4.8, 5.4, 5.6, 6.2, 7.4, 7.2, 7.2, 7.2, 7.2
}
