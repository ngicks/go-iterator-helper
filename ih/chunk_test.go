package ih_test

import (
	"iter"
	"testing"

	"github.com/ngicks/go-iterator-helper/ih"
)

func TestChunk(t *testing.T) {
	src := []int{28, 6, 49, 65, 30, 3, 9, 5, 1}

	t.Run("divide 9 by 3", func(t *testing.T) {
		testCase1[[]int]{
			Seq: func() iter.Seq[[]int] {
				return ih.Chunk[[]int](src, 3)
			},
			Expected: [][]int{{28, 6, 49}, {65, 30, 3}, {9, 5, 1}},
			BreakAt:  2,
		}.Test(t)
	})

	t.Run("divide 9 by 4", func(t *testing.T) {
		testCase1[[]int]{
			Seq: func() iter.Seq[[]int] {
				return ih.Chunk[[]int](src, 4)
			},
			Expected: [][]int{{28, 6, 49, 65}, {30, 3, 9, 5}, {1}},
			BreakAt:  2,
		}.Test(t)
	})

	t.Run("divide 9 by 0", func(t *testing.T) {
		testCase1[[]int]{
			Seq: func() iter.Seq[[]int] {
				return ih.Chunk[[]int](src, 0)
			},
			Expected: nil,
		}.Test(t)
	})

	t.Run("divide 9 by -1", func(t *testing.T) {
		testCase1[[]int]{
			Seq: func() iter.Seq[[]int] {
				return ih.Chunk[[]int](src, -1)
			},
			Expected: nil,
		}.Test(t)
	})
}
