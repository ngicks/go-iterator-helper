package goiteratorhelper_test

import (
	"fmt"

	"github.com/ngicks/go-iterator-helper/hiter"
	"github.com/ngicks/go-iterator-helper/hiter/iterable"
)

func Example_peek_and_continue() {
	// iterator that yields 0 to 9 sequentially.
	src := hiter.Range(0, 10)

	fmt.Println("It replays data if break-ed and resumed.")

	count := 3
	first := true
	for v := range src {
		count--
		if count < 0 {
			break
		}
		if !first {
			fmt.Print(", ")
		}
		first = false
		fmt.Printf("%d", v)
	}
	fmt.Println()
	fmt.Println("break and resume")
	first = true
	for v := range hiter.Limit(3, src) {
		if !first {
			fmt.Print(", ")
		}
		first = false
		fmt.Printf("%d", v)
	}
	fmt.Print("\n\n")

	fmt.Println("converting it to be resumable.")
	resumable := iterable.NewResumable(src)

	v0, _ := hiter.First(resumable.IntoIter())
	fmt.Printf("first:  %d\n", v0)
	v1, _ := hiter.First(resumable.IntoIter())
	fmt.Printf("second: %d\n", v1)

	fmt.Println()
	fmt.Println("reconnect them to whole iterator.")
	first = true
	for v := range hiter.Concat(hiter.Once(v0), hiter.Once(v1), resumable.IntoIter()) {
		if !first {
			fmt.Print(", ")
		}
		first = false
		fmt.Printf("%d", v)
	}
	fmt.Println()

	fmt.Println("\nYou can achieve above also with iterable.Peekable")
	peekable := iterable.NewPeekable(src)
	fmt.Printf("%#v\n", peekable.Peek(5))
	first = true
	for v := range peekable.IntoIter() {
		if !first {
			fmt.Print(", ")
		}
		first = false
		fmt.Printf("%d", v)
	}
	fmt.Println()

	// Output:
	// It replays data if break-ed and resumed.
	// 0, 1, 2
	// break and resume
	// 0, 1, 2
	//
	// converting it to be resumable.
	// first:  0
	// second: 1
	//
	// reconnect them to whole iterator.
	// 0, 1, 2, 3, 4, 5, 6, 7, 8, 9
	//
	// You can achieve above also with iterable.Peekable
	// []int{0, 1, 2, 3, 4}
	// 0, 1, 2, 3, 4, 5, 6, 7, 8, 9
}
