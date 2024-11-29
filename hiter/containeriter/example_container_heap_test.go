package containeriter_test

import (
	"cmp"
	"container/heap"
	"fmt"

	"github.com/ngicks/go-iterator-helper/hiter/containeriter"
	"github.com/ngicks/go-iterator-helper/x/exp/xiter"
)

var _ heap.Interface = (*sliceHeap)(nil)

type sliceHeap []int

func (s *sliceHeap) Len() int           { return len(*s) }
func (s *sliceHeap) Less(i, j int) bool { return cmp.Less((*s)[i], (*s)[j]) }
func (s *sliceHeap) Swap(i, j int)      { (*s)[i], (*s)[j] = (*s)[j], (*s)[i] }
func (s *sliceHeap) Push(x any)         { (*s) = append((*s), x.(int)) }
func (s *sliceHeap) Pop() any {
	p := (*s)[len(*s)-1]
	// Zeroing out the removed part of slice.
	// This does nothing for types like int.
	// If the type is pointer or struct that contains pointer,
	// zeroing out lets GC clean up now-unused elements.
	(*s)[len(*s)-1] = 0
	*s = (*s)[:len(*s)-1]
	return p
}

func ExampleHeap() {
	h := &sliceHeap{0, 6, 1, 3, 2, 8, 210, 3, 7, 9, 2, 1, 54, 7}
	heap.Init(h)

	for num := range xiter.Limit(containeriter.Heap[int](h), 5) {
		fmt.Printf("%d, ", num)
	}

	fmt.Println("...stopped here. and...")

	for num := range containeriter.Heap[int](h) {
		fmt.Printf("%d, ", num)
		if h.Len() == 1 {
			break
		}
	}

	for num := range containeriter.Heap[int](h) {
		fmt.Printf("%d\n", num)
	}

	// Output:
	// 0, 1, 1, 2, 2, ...stopped here. and...
	// 3, 3, 6, 7, 7, 8, 9, 54, 210
}
