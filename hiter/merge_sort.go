package hiter

import (
	"cmp"
	"fmt"
	"iter"
	"slices"

	"github.com/ngicks/go-iterator-helper/x/exp/xiter"
)

// MergeSort implements merge sort algorithm.
// Basically you should use []T -> []T implementation since this allocates a lot more.
// MergeSort is worthy only when T is big structs and you only need one element at time, not a whole slice.
func MergeSort[S ~[]T, T cmp.Ordered](m S) iter.Seq[T] {
	return MergeSortFunc(m, cmp.Compare)
}

// MergeSortFunc is like [MergeSort] but uses comparison function.
func MergeSortFunc[S ~[]T, T any](m S, cmp func(l, r T) int) iter.Seq[T] {
	if len(m) <= 1 {
		return slices.Values(m)
	}
	// TODO: We might want to implement non-recursive version? AFAIK Go does not optimize for recursion.
	return func(yield func(T) bool) {
		left, right := m[:len(m)/2], m[len(m)/2:]
		leftIter := MergeSortFunc(left, cmp)
		rightIter := MergeSortFunc(right, cmp)
		for t := range xiter.MergeFunc(leftIter, rightIter, cmp) {
			if !yield(t) {
				return
			}
		}
	}
}

// MergeSortSliceLike is like [MergeSort] that uses [SliceLike] interface instead of []T.
// This implementation is quite experimental.
// Basically you do not want to use this since it is much, much less performant.
func MergeSortSliceLike[S SliceLike[T], T cmp.Ordered](s S) iter.Seq[T] {
	return MergeSortSliceLikeFunc(s, cmp.Compare)
}

// MergeSortSliceLikeFunc is like [MergeSortSliceLike] but uses comparison function instead.
func MergeSortSliceLikeFunc[S SliceLike[T], T any](s S, cmp func(l, r T) int) iter.Seq[T] {
	return mergeSortSubbableFunc(newSubbable(s), cmp)
}

type SliceLike[T any] interface {
	Atter[T]
	Len() int
}

// ConcatSliceLike logically concatenates [SliceLike][T] implementations to form a single concatenated [SliceLike][T].
func ConcatSliceLike[T any](s ...SliceLike[T]) SliceLike[T] {
	var (
		ss          []indexRanged[T]
		totalLength int
	)

	for _, s := range s {
		ss = append(ss, indexRanged[T]{
			Start: totalLength,
			End:   totalLength + s.Len(),
			S:     s,
		})
		totalLength += s.Len()
	}

	return combiner[T]{ss, totalLength}
}

type combiner[T any] struct {
	ss          []indexRanged[T]
	totalLength int
}

type indexRanged[T any] struct {
	Start, End int
	S          SliceLike[T]
}

func compareRange[T any](r indexRanged[T], t int) int {
	switch {
	case t < r.Start:
		return 1
	case r.Start <= t && t < r.End:
		return 0
	default: // r.End <= off:
		return -1
	}
}

func (c combiner[T]) At(i int) T {
	// I took benchmarks and noticed that if ss is holding less than like 20 elements,
	// linear search is slightly faster(nano sec order).
	// At least for Go 1.23.0.
	// But I don't care about that here.
	found, ok := slices.BinarySearchFunc(c.ss, i, compareRange)
	if !ok {
		panic(fmt.Sprintf("index out of range [%d], with length of %d", i, c.totalLength))
	}
	target := c.ss[found]
	return target.S.At(i - target.Start)
}

func (c combiner[T]) Len() int {
	return c.totalLength
}

func mergeSortSubbableFunc[S SliceLike[T], T any](s subbable[S, T], cmp func(l, r T) int) iter.Seq[T] {
	if s.Len() <= 1 {
		return OmitF(IndexAccessible(s, Range(0, s.Len())))
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

type subbable[S SliceLike[T], T any] struct {
	S    S
	i, j int
}

func newSubbable[S SliceLike[T], T any](s S) subbable[S, T] {
	return subbable[S, T]{
		S: s,
		i: 0,
		j: s.Len(),
	}
}

func (s subbable[S, T]) At(i int) T {
	i = s.i + i
	if i < s.i || i >= s.j {
		panic(fmt.Sprintf("index out of range [%d], with range of [%d, %d)", i, s.i, s.j))
	}
	return s.S.At(i)
}

func (s subbable[S, T]) Len() int {
	return s.j - s.i
}

func (s subbable[S, T]) Sub(i, j int) subbable[S, T] {
	i = i + s.i
	j = j + s.i
	if i < s.i || j > s.j || i > j {
		panic(fmt.Errorf("slice bounds out of range [%d:%d] with range [%d, %d)", i, j, s.i, s.j))
	}
	return subbable[S, T]{
		S: s.S,
		i: i,
		j: j,
	}
}
