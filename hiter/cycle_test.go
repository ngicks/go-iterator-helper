package hiter_test

import (
	"iter"
	"testing"

	"github.com/ngicks/go-iterator-helper/hiter"
	"github.com/ngicks/go-iterator-helper/hiter/internal/testcase"
	"github.com/ngicks/go-iterator-helper/hiter/iterable"
)

func TestCycle(t *testing.T) {
	testcase.One[int]{
		Seq: func() iter.Seq[int] {
			return hiter.Limit(15, hiter.Cycle(hiter.Range(0, 3)))
		},
		Expected: []int{0, 1, 2, 0, 1, 2, 0, 1, 2, 0, 1, 2, 0, 1, 2},
		BreakAt:  7,
	}.Test(t)
	testcase.One[int]{
		Seq: func() iter.Seq[int] {
			return hiter.Limit(
				15,
				hiter.CycleBuffered(
					iterable.NewResumable(
						hiter.Range(0, 3),
					).IntoIter(),
				),
			)
		},
		Expected: []int{0, 1, 2, 0, 1, 2, 0, 1, 2, 0, 1, 2, 0, 1, 2},
		BreakAt:  7,
		Stateful: true,
	}.Test(t)
}

func TestCycle2(t *testing.T) {
	expected := []hiter.KeyValue[int, int]{
		{K: 0, V: 0},
		{K: 1, V: 1},
		{K: 2, V: 2},
		{K: 0, V: 0},
		{K: 1, V: 1},
		{K: 2, V: 2},
		{K: 0, V: 0},
		{K: 1, V: 1},
		{K: 2, V: 2},
		{K: 0, V: 0},
		{K: 1, V: 1},
		{K: 2, V: 2},
		{K: 0, V: 0},
		{K: 1, V: 1},
		{K: 2, V: 2},
	}

	testcase.Two[int, int]{
		Seq: func() iter.Seq2[int, int] {
			return hiter.Limit2(15, hiter.Cycle2(hiter.Enumerate(hiter.Range(0, 3))))
		},
		Seqs: []func() iter.Seq2[int, int]{
			func() iter.Seq2[int, int] {
				return hiter.Limit2(
					15,
					hiter.CycleBuffered2(
						hiter.Enumerate(hiter.Range(0, 3)),
					),
				)
			},
		},
		Expected: expected,
		BreakAt:  7,
	}.Test(t)

	testcase.Two[int, int]{
		Seq: func() iter.Seq2[int, int] {
			return hiter.Limit2(
				15,
				hiter.CycleBuffered2(
					iterable.NewResumable2(
						hiter.Enumerate(hiter.Range(0, 3)),
					).IntoIter2(),
				),
			)
		},
		Expected: expected,
		BreakAt:  7,
		Stateful: true,
	}.Test(t)
}
