package hiter_test

import (
	"cmp"
	"container/heap"
	"container/list"
	"container/ring"
	"iter"
	"slices"
	"testing"

	"github.com/ngicks/go-iterator-helper/hiter"
	"github.com/ngicks/go-iterator-helper/hiter/iterable"
)

var _ heap.Interface = (*sliceHeap)(nil)

type sliceHeap []int

func (s *sliceHeap) Len() int           { return len(*s) }
func (s *sliceHeap) Less(i, j int) bool { return cmp.Less((*s)[i], (*s)[j]) }
func (s *sliceHeap) Swap(i, j int)      { (*s)[i], (*s)[j] = (*s)[j], (*s)[i] }
func (s *sliceHeap) Push(x any)         { (*s) = append((*s), x.(int)) }
func (s *sliceHeap) Pop() any {
	p := (*s)[len(*s)-1]
	*s = slices.Delete(*s, len(*s)-1, len(*s))
	return p
}

func TestContainerHeap(t *testing.T) {
	h := &sliceHeap{2, 7, 9, 0, 9, 1}
	heap.Init(h)
	testCase1[int]{
		Seq: func() iter.Seq[int] {
			h := slices.Clone(*h)
			return hiter.Heap[int](&h)
		},
		Seqs: []func() iter.Seq[int]{
			func() iter.Seq[int] {
				h := slices.Clone(*h)
				return iterable.Heap[int]{Interface: &h}.IntoIter()
			},
		},
		Expected: []int{0, 1, 2, 7, 9, 9},
		BreakAt:  3,
		Stateful: true,
	}.Test(t)
}

func TestContainerList(t *testing.T) {
	s := list.New()
	for i := range 5 {
		s.PushBack(i + 5)
	}

	t.Run("ListAll", func(t *testing.T) {
		testCase1[int]{
			Seq: func() iter.Seq[int] {
				return hiter.ListAll[int](s)
			},
			Seqs: []func() iter.Seq[int]{
				func() iter.Seq[int] {
					return iterable.ListAll[int]{List: s}.Iter()
				},
			},
			Expected: []int{5, 6, 7, 8, 9},
			BreakAt:  3,
		}.Test(t)
	})

	t.Run("ListElementAll", func(t *testing.T) {
		testCase1[int]{
			Seq: func() iter.Seq[int] {
				return hiter.ListElementAll[int](s.Front().Next().Next())
			},
			Seqs: []func() iter.Seq[int]{
				func() iter.Seq[int] {
					return iterable.ListElementAll[int]{Element: s.Front().Next().Next()}.Iter()
				},
			},
			Expected: []int{7, 8, 9},
			BreakAt:  2,
		}.Test(t)
	})

	t.Run("ListBackward", func(t *testing.T) {
		testCase1[int]{
			Seq: func() iter.Seq[int] {
				return hiter.ListBackward[int](s)
			},
			Seqs: []func() iter.Seq[int]{
				func() iter.Seq[int] {
					return iterable.ListBackward[int]{List: s}.Iter()
				},
			},
			Expected: []int{9, 8, 7, 6, 5},
			BreakAt:  3,
		}.Test(t)
	})

	t.Run("ListElementBackward", func(t *testing.T) {
		testCase1[int]{
			Seq: func() iter.Seq[int] {
				return hiter.ListElementBackward[int](s.Back().Prev())
			},
			Seqs: []func() iter.Seq[int]{
				func() iter.Seq[int] {
					return iterable.ListElementBackward[int]{Element: s.Back().Prev()}.Iter()
				},
			},
			Expected: []int{8, 7, 6, 5},
			BreakAt:  3,
		}.Test(t)
	})
}

func TestContainerRing(t *testing.T) {
	r := ring.New(5)
	r.Value = 5
	i := 6
	for n := r.Next(); ; n = n.Next() {
		if n == r {
			break
		}
		n.Value = i
		i++
	}

	t.Run("RingAll", func(t *testing.T) {
		testCase1[int]{
			Seq: func() iter.Seq[int] {
				return hiter.RingAll[int](r.Move(2))
			},
			Seqs: []func() iter.Seq[int]{
				func() iter.Seq[int] {
					return iterable.Ring[int]{Ring: r.Move(2)}.Iter()
				},
			},
			Expected: []int{7, 8, 9, 5, 6},
			BreakAt:  3,
		}.Test(t)
	})

	t.Run("RingBackward", func(t *testing.T) {
		testCase1[int]{
			Seq: func() iter.Seq[int] {
				return hiter.RingBackward[int](r.Move(2))
			},
			Seqs: []func() iter.Seq[int]{
				func() iter.Seq[int] {
					return iterable.RingBackward[int]{Ring: r.Move(2)}.Iter()
				},
			},
			Expected: []int{7, 6, 5, 9, 8},
			BreakAt:  3,
		}.Test(t)
	})
}
