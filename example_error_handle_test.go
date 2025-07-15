package goiteratorhelper_test

import (
	"errors"
	"fmt"
	"slices"

	"github.com/ngicks/go-iterator-helper/hiter"
	"github.com/ngicks/go-iterator-helper/hiter/errbox"
	"github.com/ngicks/go-iterator-helper/hiter/mapper"
)

// Example error handle demonstrates various way to handle error.
func Example_error_handle() {
	var (
		errSample  = errors.New("sample")
		errSample2 = errors.New("sample2")
	)

	erroneous := hiter.Pairs(
		hiter.Range(0, 6),
		hiter.Concat(
			hiter.Repeat(error(nil), 2),
			hiter.Repeat(errSample2, 2),
			hiter.Once(errSample),
			hiter.Once(error(nil)),
		),
	)

	fmt.Println("TryFind:")
	v, idx, err := hiter.TryFind(func(i int) bool { return i > 0 }, erroneous)
	fmt.Printf("v = %d, idx = %d, err = %v\n", v, idx, err)
	v, idx, err = hiter.TryFind(func(i int) bool { return i > 5 }, erroneous)
	fmt.Printf("v = %d, idx = %d, err = %v\n", v, idx, err)
	fmt.Println()

	fmt.Println("TryForEach:")
	err = hiter.TryForEach(func(i int) { fmt.Printf("i = %d\n", i) }, erroneous)
	fmt.Printf("err = %v\n", err)
	fmt.Println()

	fmt.Println("TryReduce:")
	collected, err := hiter.TryReduce(func(c []int, i int) []int { return append(c, i) }, nil, erroneous)
	fmt.Printf("collected = %#v, err = %v\n", collected, err)
	fmt.Println()

	fmt.Println("HandleErr:")
	var handled error
	collected = slices.Collect(
		mapper.HandleErr(
			func(i int, err error) bool {
				handled = err
				return errors.Is(err, errSample2)
			},
			erroneous,
		),
	)
	fmt.Printf("collected = %#v, err = %v\n", collected, handled)
	fmt.Println()

	fmt.Println("*errbox.Box:")
	box := errbox.New(erroneous)
	collected = slices.Collect(box.IntoIter())
	fmt.Printf("collected = %#v, err = %v\n", collected, box.Err())
	fmt.Println()
	// Output:
	// TryFind:
	// v = 1, idx = 1, err = <nil>
	// v = 0, idx = -1, err = sample2
	//
	// TryForEach:
	// i = 0
	// i = 1
	// err = sample2
	//
	// TryReduce:
	// collected = []int{0, 1}, err = sample2
	//
	// HandleErr:
	// collected = []int{0, 1}, err = sample
	//
	// *errbox.Box:
	// collected = []int{0, 1}, err = sample2
}
