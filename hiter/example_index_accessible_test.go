package hiter_test

import (
	"fmt"

	"github.com/ngicks/go-iterator-helper/hiter"
)

type atSliceStr []string

func (a atSliceStr) At(i int) string {
	return a[i]
}

var atSliceSrc = atSliceStr{
	"foo", "bar", "baz",
	"qux", "quux", "corge",
	"grault", "garply", "waldo",
	"fred", "plugh", "xyzzy",
	"thud",
}

func ExampleIndexAccessible() {
	for i, s := range hiter.IndexAccessible(atSliceSrc, hiter.Range(3, 7)) {
		fmt.Printf("%d: %s\n", i, s)
	}
	// Output:
	// 3: qux
	// 4: quux
	// 5: corge
	// 6: grault
}
