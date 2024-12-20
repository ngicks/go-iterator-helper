package containeriter_test

import (
	"container/heap"
	"container/list"
	"container/ring"
	"iter"
	"slices"
	"testing"

	"github.com/ngicks/go-iterator-helper/hiter/containeriter"
	"github.com/ngicks/go-iterator-helper/hiter/internal/testcase"
	"github.com/ngicks/go-iterator-helper/hiter/iterable"
)

func TestContainerHeap(t *testing.T) {
	h := &sliceHeap{2, 7, 9, 0, 9, 1}
	heap.Init(h)
	testcase.One[int]{
		Seq: func() iter.Seq[int] {
			h := slices.Clone(*h)
			return containeriter.Heap[int](&h)
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
		testcase.One[int]{
			Seq: func() iter.Seq[int] {
				return containeriter.ListAll[int](s)
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
		testcase.One[int]{
			Seq: func() iter.Seq[int] {
				return containeriter.ListElementAll[int](s.Front().Next().Next())
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
		testcase.One[int]{
			Seq: func() iter.Seq[int] {
				return containeriter.ListBackward[int](s)
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
		testcase.One[int]{
			Seq: func() iter.Seq[int] {
				return containeriter.ListElementBackward[int](s.Back().Prev())
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
		testcase.One[int]{
			Seq: func() iter.Seq[int] {
				return containeriter.RingAll[int](r.Move(2))
			},
			Seqs: []func() iter.Seq[int]{
				func() iter.Seq[int] {
					return iterable.RingAll[int]{Ring: r.Move(2)}.Iter()
				},
			},
			Expected: []int{7, 8, 9, 5, 6},
			BreakAt:  3,
		}.Test(t)
	})

	t.Run("RingBackward", func(t *testing.T) {
		testcase.One[int]{
			Seq: func() iter.Seq[int] {
				return containeriter.RingBackward[int](r.Move(2))
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
