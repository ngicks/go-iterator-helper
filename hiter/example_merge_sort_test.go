package hiter_test

import (
	"cmp"
	"fmt"
	"iter"
	"math/rand/v2"
	"slices"

	"github.com/ngicks/go-iterator-helper/hiter"
	"github.com/ngicks/go-iterator-helper/x/exp/xiter"
)

// avoiding xter dep
func limit[V any](seq iter.Seq[V], n int) iter.Seq[V] {
	return func(yield func(V) bool) {
		if n <= 0 {
			return
		}
		for v := range seq {
			if !yield(v) {
				return
			}
			if n--; n <= 0 {
				break
			}
		}
	}
}

func mergeSortFunc[S ~[]T, T any](m S, cmp func(l, r T) int) S {
	if len(m) <= 1 {
		return m
	}
	left, right := m[:len(m)/2], m[len(m)/2:]
	left = mergeSortFunc(left, cmp)
	right = mergeSortFunc(right, cmp)
	return mergeFunc(left, right, cmp)
}

func mergeFunc[S ~[]T, T any](l, r S, cmp func(l, r T) int) S {
	m := make(S, len(l)+len(r))
	var i int
	for i = 0; len(l) > 0 && len(r) > 0; i++ {
		if cmp(l[0], r[0]) < 0 {
			m[i] = l[0]
			l = l[1:]
		} else {
			m[i] = r[0]
			r = r[1:]
		}
	}
	for _, t := range l {
		m[i] = t
		i++
	}
	for _, t := range r {
		m[i] = t
		i++
	}
	return m
}

func mergeSortIterFunc[S ~[]T, T any](m S, cmp func(l, r T) int) iter.Seq[T] {
	if len(m) <= 1 {
		return slices.Values(m)
	}
	return func(yield func(T) bool) {
		left, right := m[:len(m)/2], m[len(m)/2:]
		leftIter := mergeSortIterFunc(left, cmp)
		rightIter := mergeSortIterFunc(right, cmp)
		for t := range xiter.MergeFunc(leftIter, rightIter, cmp) {
			if !yield(t) {
				return
			}
		}
	}
}

type SliceLike[T any] interface {
	hiter.Atter[T]
	Len() int
}

var _ SliceLike[any] = sliceAdapter[any]{}

type sliceAdapter[T any] []T

func (s sliceAdapter[T]) At(i int) T {
	return s[i]
}

func (s sliceAdapter[T]) Len() int {
	return len(s)
}

type subbable[S SliceLike[T], T any] struct {
	S    S
	i, j int
}

func (s subbable[S, T]) At(i int) T {
	i = s.i + i
	if i >= s.j {
		panic("index out of range")
	}
	return s.S.At(i)
}

func (s subbable[S, T]) Len() int {
	return s.j - s.i
}

func (s subbable[S, T]) Sub(i, j int) subbable[S, T] {
	i = i + s.i
	j = j + s.i
	if i < 0 || j > s.j || i > j {
		panic(fmt.Errorf("index out of range: i=%d, j=%d,len=%d", i, j, s.Len()))
	}
	return subbable[S, T]{
		S: s.S,
		i: i,
		j: j,
	}
}

func mergeSortAtterFunc[S SliceLike[T], T any](s S, cmp func(l, r T) int) iter.Seq[T] {
	sub := subbable[S, T]{
		S: s,
		i: 0,
		j: s.Len(),
	}
	return mergeSortSubbableFunc(sub, cmp)
}

func mergeSortSubbableFunc[S SliceLike[T], T any](s subbable[S, T], cmp func(l, r T) int) iter.Seq[T] {
	if s.Len() <= 1 {
		return hiter.OmitF(hiter.IndexAccessible(s, hiter.Range(0, s.Len())))
	}
	return func(yield func(T) bool) {
		left, right := s.Sub(0, s.Len()/2), s.Sub(s.Len()/2, s.Len())
		leftIter := mergeSortSubbableFunc(left, cmp)
		rightIter := mergeSortSubbableFunc(right, cmp)
		for t := range xiter.MergeFunc(leftIter, rightIter, cmp) {
			if !yield(t) {
				return
			}
		}
	}
}

// Example_merge_sort implements slice version merge sort and re-implements iterator version of it.
func Example_merge_sort() {
	rng := hiter.RepeatFunc(func() int { return rand.N(20) }, -1)
	fmt.Printf("merge sort: %t\n",
		slices.IsSorted(mergeSortFunc(slices.Collect(limit(rng, 10)), cmp.Compare)),
	)
	fmt.Printf(
		"merge sort iter: %t\n",
		slices.IsSorted(
			slices.Collect(
				mergeSortIterFunc(slices.Collect(limit(rng, 10)), cmp.Compare),
			),
		),
	)
	fmt.Printf(
		"merge sort atter: %t\n",
		slices.IsSorted(
			slices.Collect(
				mergeSortAtterFunc(sliceAdapter[int](slices.Collect(limit(rng, 10))), cmp.Compare),
			),
		),
	)
	// Output:
	// merge sort: true
	// merge sort iter: true
	// merge sort atter: true
}
