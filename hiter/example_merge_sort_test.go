package hiter_test

import (
	"cmp"
	"fmt"
	"iter"
	"slices"

	"github.com/ngicks/go-iterator-helper/hiter"
	"github.com/ngicks/go-iterator-helper/hiter/sh"
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

var _ hiter.SliceLike[any] = sliceAdapter[any]{}

type sliceAdapter[T any] []T

func (s sliceAdapter[T]) At(i int) T {
	return s[i]
}

func (s sliceAdapter[T]) Len() int {
	return len(s)
}

func ExampleMergeSort() {
	rng := sh.Rng(20)
	fmt.Printf("merge sort: %t\n",
		slices.IsSorted(mergeSortFunc(slices.Collect(limit(rng, 10)), cmp.Compare)),
	)
	fmt.Printf(
		"merge sort iter: %t\n",
		slices.IsSorted(
			slices.Collect(
				hiter.MergeSort(slices.Collect(limit(rng, 10))),
			),
		),
	)
	fmt.Printf(
		"merge sort atter: %t\n",
		slices.IsSorted(
			slices.Collect(
				hiter.MergeSortSliceLike(sliceAdapter[int](slices.Collect(limit(rng, 10)))),
			),
		),
	)
	fmt.Printf(
		"merge sort atter logically concatenated: %t\n",
		slices.IsSorted(
			slices.Collect(
				hiter.MergeSortSliceLike(
					hiter.ConcatSliceLike(
						slices.Collect(
							xiter.Map(
								func(i int) hiter.SliceLike[int] {
									return sliceAdapter[int](slices.Collect(limit(rng, i)))
								},
								xiter.Limit(rng, 5),
							),
						)...,
					),
				),
			),
		),
	)
	// Output:
	// merge sort: true
	// merge sort iter: true
	// merge sort atter: true
	// merge sort atter logically concatenated: true
}
