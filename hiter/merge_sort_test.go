package hiter_test

import (
	"slices"
	"testing"

	"github.com/ngicks/go-iterator-helper/hiter"
	"gotest.tools/v3/assert"
	"gotest.tools/v3/assert/cmp"
)

func TestMergeSort(t *testing.T) {
	type testCases struct {
		input    []int
		expected []int
	}
	for _, tc := range []testCases{
		{},
		{[]int{5, 6, 2, 3, 31, 7, 9, 4}, []int{2, 3, 4, 5, 6, 7, 9, 31}},
	} {
		collected := slices.Collect(hiter.MergeSort(tc.input))
		assert.Assert(t, slices.IsSorted(collected))
		assert.Assert(t, cmp.DeepEqual(tc.expected, collected))
	}
}

func TestMergeSliceLike(t *testing.T) {
	type testCases struct {
		input    hiter.SliceLike[int]
		expected []int
	}
	for _, tc := range []testCases{
		{sliceAdapter[int]([]int{}), nil},
		{sliceAdapter[int]([]int{5, 6, 2, 3, 31, 7, 9, 4}), []int{2, 3, 4, 5, 6, 7, 9, 31}},
		{
			hiter.ConcatSliceLike(
				sliceAdapter[int]([]int{5}),
				sliceAdapter[int]([]int{9, 2, 8, 3, 1}),
				sliceAdapter[int]([]int{4, 7, 6}),
			),
			[]int{1, 2, 3, 4, 5, 6, 7, 8, 9},
		},
		{
			hiter.ConcatSliceLike[int](),
			nil,
		},
		{
			hiter.ConcatSliceLike(
				sliceAdapter[int]([]int{5, 9, 2, 8, 3, 1}),
				sliceAdapter[int]([]int{}),
				sliceAdapter[int]([]int{4, 7, 6}),
			),
			[]int{1, 2, 3, 4, 5, 6, 7, 8, 9},
		},
	} {
		collected := slices.Collect(hiter.MergeSortSliceLike(tc.input))
		assert.Assert(t, slices.IsSorted(collected))
		assert.Assert(t, cmp.DeepEqual(tc.expected, collected))
	}
}
