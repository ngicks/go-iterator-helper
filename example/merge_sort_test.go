package example

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

func Example() {
	rng := hiter.RepeatFunc(func() int { return rand.N(20) }, -1)
	fmt.Printf("merged sort: %t\n", slices.IsSorted(mergeSortFunc(slices.Collect(limit(rng, 10)), cmp.Compare)))

	fmt.Printf(
		"merge sort iter: %t\n",
		slices.IsSorted(
			slices.Collect(
				mergeSortIterFunc(slices.Collect(limit(rng, 10)), cmp.Compare),
			),
		),
	)

	// Output:
	// merged sort: true
	// merge sort iter: true
}
