package hiter_test

import (
	"iter"
	"slices"
	"testing"

	"github.com/ngicks/go-iterator-helper/hiter"
	"github.com/ngicks/go-iterator-helper/hiter/internal/testcase"
	"github.com/ngicks/go-iterator-helper/x/exp/xiter"
)

func TestUniqueFunc(t *testing.T) {
	testcase.One[[]int]{
		Seq: func() iter.Seq[[]int] {
			return hiter.UniqueFunc(
				func(v []int) int {
					return hiter.Sum(slices.Values(v))
				},
				slices.Values([][]int{
					{1, 2},
					{5},
					{2, 1},
					{4},
					{3, 0},
				}),
			)
		},
		BreakAt:  2,
		Expected: [][]int{{1, 2}, {5}, {4}},
	}.Test(t)
}

func TestUnique(t *testing.T) {
	testcase.One[int]{
		Seq: func() iter.Seq[int] {
			return hiter.Unique(slices.Values([]int{1, 2, 3, 3, 5, 1, 2, 5, 4}))
		},
		BreakAt:  2,
		Expected: []int{1, 2, 3, 5, 4},
	}.Test(t)
}

func TestUniqueFunc2(t *testing.T) {
	testcase.Two[[]int, []int]{
		Seq: func() iter.Seq2[[]int, []int] {
			return hiter.UniqueFunc2(
				func(k, v []int) int {
					return hiter.Sum(xiter.Concat(slices.Values(k), slices.Values(v)))
				},
				hiter.Values2([]hiter.KeyValue[[]int, []int]{
					{[]int{1, 2}, []int{3}},
					{[]int{6}, []int{0}},
					{[]int{1}, []int{1}},
					{[]int{2}, []int{2}},
					{[]int{3}, []int{3}},
				}),
			)
		},
		BreakAt: 2,
		Expected: []hiter.KeyValue[[]int, []int]{
			{[]int{1, 2}, []int{3}},
			{[]int{1}, []int{1}},
			{[]int{2}, []int{2}},
		},
	}.Test(t)
}

func TestUnique2(t *testing.T) {
	testcase.Two[int, int]{
		Seq: func() iter.Seq2[int, int] {
			return hiter.Unique2(
				hiter.Values2(
					[]hiter.KeyValue[int, int]{
						{K: 1, V: 0},
						{K: 2, V: 0},
						{K: 3, V: 0},
						{K: 3, V: 0},
						{K: 5, V: 0},
						{K: 1, V: 1},
						{K: 2, V: 0},
						{K: 5, V: 0},
						{K: 4, V: 0},
					},
				),
			)
		},
		BreakAt: 2,
		Expected: []hiter.KeyValue[int, int]{
			{K: 1, V: 0},
			{K: 2, V: 0},
			{K: 3, V: 0},
			{K: 5, V: 0},
			{K: 1, V: 1},
			{K: 4, V: 0},
		},
	}.Test(t)
}
