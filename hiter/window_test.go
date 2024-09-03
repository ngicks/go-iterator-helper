package hiter_test

import (
	"iter"
	"testing"

	"github.com/ngicks/go-iterator-helper/hiter"
)

func TestWindow(t *testing.T) {
	src := []int{28, 6, 49, 65, 30, 3}

	t.Run("window 9 by 3", func(t *testing.T) {
		testCase1[[]int]{
			Seq: func() iter.Seq[[]int] {
				return hiter.Window(src, 3)
			},
			Expected: [][]int{{28, 6, 49}, {6, 49, 65}, {49, 65, 30}, {65, 30, 3}},
			BreakAt:  2,
		}.Test(t)
	})

	t.Run("window 9 by 4", func(t *testing.T) {
		testCase1[[]int]{
			Seq: func() iter.Seq[[]int] {
				return hiter.Window(src, 4)
			},
			Expected: [][]int{{28, 6, 49, 65}, {6, 49, 65, 30}, {49, 65, 30, 3}},
			BreakAt:  2,
		}.Test(t)
	})

	t.Run("window 9 by 10", func(t *testing.T) {
		testCase1[[]int]{
			Seq: func() iter.Seq[[]int] {
				return hiter.Window(src, 10)
			},
			Expected: [][]int{{28, 6, 49, 65, 30, 3}},
		}.Test(t)
	})

	t.Run("window 9 by 0", func(t *testing.T) {
		testCase1[[]int]{
			Seq: func() iter.Seq[[]int] {
				return hiter.Window(src, 0)
			},
			Expected: nil,
		}.Test(t)
	})

	t.Run("window 9 by -1", func(t *testing.T) {
		testCase1[[]int]{
			Seq: func() iter.Seq[[]int] {
				return hiter.Window(src, -1)
			},
			Expected: nil,
		}.Test(t)
	})
}
