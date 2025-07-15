package goiteratorhelper

import (
	"fmt"

	"github.com/ngicks/go-iterator-helper/hiter"
	"github.com/ngicks/go-iterator-helper/hiter/stringsiter"
)

func ExampleRange_char() {
	fmt.Println(
		stringsiter.Collect(
			hiter.Map(
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
