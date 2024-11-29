package containeriter_test

import (
	"container/ring"
	"fmt"
	"iter"
	"slices"

	"github.com/ngicks/go-iterator-helper/hiter"
	"github.com/ngicks/go-iterator-helper/hiter/containeriter"
	"github.com/ngicks/go-iterator-helper/x/exp/xiter"
)

func ExampleRingAll() {
	ringBufSize := 5
	r := ring.New(ringBufSize)

	for i := range ringBufSize {
		r.Value = i
		r = r.Next()
	}

	fmt.Printf("all:      %#v\n", slices.Collect(containeriter.RingAll[int](r)))
	fmt.Printf("backward: %#v\n", slices.Collect(containeriter.RingBackward[int](r.Prev())))

	// Now, we'll demonstrate buffer like usage of ring.

	pushBack := func(v int) {
		r.Value = v
		r = r.Next()
	}

	pushBack(12)
	pushBack(5)
	fmt.Printf("1:        %#v\n", slices.Collect(containeriter.RingAll[int](r)))

	pushBack(8)
	fmt.Printf("2:        %#v\n", slices.Collect(containeriter.RingAll[int](r)))

	pushBack(36)
	fmt.Printf("3:        %#v\n", slices.Collect(containeriter.RingAll[int](r)))

	// Output:
	// all:      []int{0, 1, 2, 3, 4}
	// backward: []int{4, 3, 2, 1, 0}
	// 1:        []int{2, 3, 4, 12, 5}
	// 2:        []int{3, 4, 12, 5, 8}
	// 3:        []int{4, 12, 5, 8, 36}
}

func dynamicWindow[V any](size int, seq iter.Seq[V]) iter.Seq[iter.Seq[V]] {
	return func(yield func(iter.Seq[V]) bool) {
		var (
			full     = false
			bufStart = ring.New(size)
		)
		buf := bufStart
		for v := range seq {
			if !full {
				buf.Value = v
				buf = buf.Next()
				if buf == bufStart {
					full = true
					if !yield(containeriter.RingAll[V](buf)) {
						return
					}
				}
				continue
			}
			buf.Value = v
			buf = buf.Next()
			if !yield(containeriter.RingAll[V](buf)) {
				return
			}
		}
	}
}

// ExampleRingAll_moving_average demonstrates buffer-like usage of RingAll.
func ExampleRingAll_moving_average() {
	windowSize := 5
	src := slices.Values([]int{1, 0, 1, 0, 1, 0, 5, 3, 2, 3, 4, 6, 5, 3, 6, 7, 7, 8, 9, 5, 7, 7, 8})

	first := true
	for avg := range xiter.Map(
		func(s iter.Seq[int]) float64 {
			return float64(hiter.Sum(s)) / float64(windowSize)
		},
		dynamicWindow(windowSize, src),
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
