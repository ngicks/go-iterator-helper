package hiter_test

import (
	"fmt"
	"sync/atomic"

	"github.com/ngicks/go-iterator-helper/hiter"
	"github.com/ngicks/go-iterator-helper/hiter/iterable"
)

func Example_decorate_string() {
	src := "foo bar baz"
	var num atomic.Int32
	numListTitle := iterable.RepeatableFunc[string]{
		FnV: func() string { return fmt.Sprintf("%d. ", num.Add(1)) },
		N:   1,
	}
	m := hiter.StringsCollect(
		hiter.SkipLast(
			hiter.Decorate(
				hiter.StringsSplitFunc(src, -1, hiter.StringsCutWord),
				numListTitle,
				iterable.Repeatable[string]{V: " ", N: 1},
			),
			1,
		),
		9+((2 /*num*/ +2 /*. */ +1 /* */)*3),
	)
	fmt.Printf("%s\n", m)
	// Output:
	// 1. foo 2. bar 3. baz
}
