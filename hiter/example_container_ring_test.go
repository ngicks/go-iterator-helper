package hiter_test

import (
	"container/ring"
	"fmt"
	"slices"

	"github.com/ngicks/go-iterator-helper/hiter"
)

func ExampleRingAll() {
	ringBufSize := 5
	r := ring.New(ringBufSize)

	for i := range ringBufSize {
		r.Value = i
		r = r.Next()
	}

	fmt.Printf("all:      %#v\n", slices.Collect(hiter.RingAll[int](r)))
	fmt.Printf("backward: %#v\n", slices.Collect(hiter.RingBackward[int](r.Prev())))

	// Now, we'll demonstrate buffer like usage of ring.

	pushBack := func(v int) {
		r.Value = v
		r = r.Next()
	}

	pushBack(12)
	pushBack(5)
	fmt.Printf("1:        %#v\n", slices.Collect(hiter.RingAll[int](r)))

	pushBack(8)
	fmt.Printf("2:        %#v\n", slices.Collect(hiter.RingAll[int](r)))

	pushBack(36)
	fmt.Printf("3:        %#v\n", slices.Collect(hiter.RingAll[int](r)))

	// Output:
	// all:      []int{0, 1, 2, 3, 4}
	// backward: []int{4, 3, 2, 1, 0}
	// 1:        []int{2, 3, 4, 12, 5}
	// 2:        []int{3, 4, 12, 5, 8}
	// 3:        []int{4, 12, 5, 8, 36}
}
