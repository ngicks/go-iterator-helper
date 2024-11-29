package goiteratorhelper

import (
	"fmt"

	"github.com/ngicks/go-iterator-helper/hiter"
	"github.com/ngicks/go-iterator-helper/hiter/stringsiter"
	"github.com/ngicks/go-iterator-helper/x/exp/xiter"
)

func ExampleRange_char() {
	fmt.Println(
		stringsiter.StringsCollect(
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
