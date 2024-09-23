package hiter_test

import (
	"iter"
	"slices"
	"testing"

	"github.com/ngicks/go-iterator-helper/hiter"
	"github.com/ngicks/go-iterator-helper/x/exp/xiter"
	"gotest.tools/v3/assert"
	"gotest.tools/v3/assert/cmp"
)

type windowTestCase struct {
	name     string
	size     int
	expected [][]int
	breakAt  int
}

var (
	widowSrc        = []int{28, 6, 49, 65, 30, 3}
	windowTestCases = []windowTestCase{
		{
			name:     "window 9 by 3",
			size:     3,
			expected: [][]int{{28, 6, 49}, {6, 49, 65}, {49, 65, 30}, {65, 30, 3}},
			breakAt:  2,
		},
		{
			name:     "window 9 by 4",
			size:     4,
			expected: [][]int{{28, 6, 49, 65}, {6, 49, 65, 30}, {49, 65, 30, 3}},
			breakAt:  2,
		},
		{
			name: "window 9 by 10",
			size: 10,
		},
		{
			name: "window 9 by 0",
		},
		{
			name: "window 9 by -1",
			size: -1,
		},
	}
)

func TestWindow(t *testing.T) {
	t.Run("empty slice", func(t *testing.T) {
		c := slices.Collect(hiter.Window[[]int](nil, 5))
		assert.Assert(t, cmp.Len(c, 0))
	})

	for _, tc := range windowTestCases {
		t.Run(tc.name, func(t *testing.T) {
			testCase1[[]int]{
				Seq: func() iter.Seq[[]int] {
					return hiter.Window(widowSrc, tc.size)
				},
				Expected: tc.expected,
				BreakAt:  tc.breakAt,
			}.Test(t)
		})
	}
}

func TestWindowSeq(t *testing.T) {
	t.Run("empty slice", func(t *testing.T) {
		c := slices.Collect(xiter.Map(slices.Collect, hiter.WindowSeq(5, hiter.Repeat(0, 0))))
		assert.Assert(t, cmp.Len(c, 0))
	})

	for _, tc := range windowTestCases {
		t.Run(tc.name, func(t *testing.T) {
			testCase1[[]int]{
				Seq: func() iter.Seq[[]int] {
					return xiter.Map(
						slices.Collect,
						hiter.WindowSeq(tc.size, slices.Values(widowSrc)),
					)
				},
				Expected: tc.expected,
				BreakAt:  tc.breakAt,
			}.Test(t)
		})
	}
}
