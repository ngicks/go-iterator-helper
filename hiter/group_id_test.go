package hiter_test

import (
	"iter"
	"slices"
	"testing"

	"github.com/ngicks/go-iterator-helper/hiter"
	"github.com/ngicks/go-iterator-helper/hiter/internal/testcase"
)

func TestWithGroupId(t *testing.T) {
	seq := slices.Values([]int{-1, 2, 5, 7, 4, -1, 6})
	testcase.Two[int, int]{
		Seq: func() iter.Seq2[int, int] {
			return hiter.WithGroupId(
				func(i int) string {
					switch {
					case i < 0:
						return "neg"
					case i%2 == 0:
						return "even"
					default:
						return "odd"
					}
				},
				seq,
			)
		},
		Expected: []hiter.KeyValue[int, int]{
			{K: 0, V: -1},
			{K: 1, V: 2},
			{K: 2, V: 5},
			{K: 2, V: 7},
			{K: 1, V: 4},
			{K: 0, V: -1},
			{K: 1, V: 6},
		},
		BreakAt: 2,
	}.Test(t)
}
