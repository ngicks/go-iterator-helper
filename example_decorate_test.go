package goiteratorhelper_test

import (
	"fmt"
	"sync/atomic"

	"github.com/ngicks/go-iterator-helper/hiter"
	"github.com/ngicks/go-iterator-helper/hiter/iterable"
	"github.com/ngicks/go-iterator-helper/hiter/stringsiter"
)

func ExampleDecorate() {
	src := "foo bar baz"
	var num atomic.Int32
	numListTitle := iterable.RepeatableFunc[string]{
		FnV: func() string { return fmt.Sprintf("%d. ", num.Add(1)) },
		N:   1,
	}
	m := stringsiter.Collect(
		hiter.SkipLast(
			1,
			hiter.Decorate(
				numListTitle,
				hiter.WrapSeqIterable(hiter.Once(" ")),
				stringsiter.SplitFunc(src, -1, stringsiter.CutWord),
			),
		),
	)
	fmt.Printf("%s\n", m)
	// Output:
	// 1. foo 2. bar 3. baz
}
