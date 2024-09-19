package hiter_test

import (
	"fmt"

	"github.com/ngicks/go-iterator-helper/hiter"
	"github.com/ngicks/go-iterator-helper/x/exp/xiter"
)

func ExampleRange_prevent_off_by_one() {
	for i := range hiter.LimitUntil(
		func(i int) bool { return i < 50 },
		xiter.Map(
			func(i int) int { return i * 7 },
			hiter.Range(0, 10),
		),
	) {
		if i > 0 {
			fmt.Printf(" ")
		}
		fmt.Printf("%d", i)
	}
	// Output:
	// 0 7 14 21 28 35 42 49
}

func ExampleRange_char() {
	fmt.Println(
		hiter.StringsCollect(
			27,
			xiter.Map(
				func(r rune) string {
					return string(r - ('a' - 'A'))
				},
				hiter.Range('a', 'z'+1),
			),
		),
	)
	// Output:
	// ABCDEFGHIJKLMNOPQRSTUVWXYZ
}
